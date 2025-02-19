import {
  Action,
  type ActionCacheClient,
  Command,
  type Digest,
  Directory,
  ExecuteResponse,
  RequestMetadata,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { AuthenticationMetadata } from "@/lib/grpc-client/buildbarn/auth/auth";
import { HistoricalExecuteResponse } from "@/lib/grpc-client/buildbarn/cas/cas";
import type {
  FileSystemAccessCacheClient,
  FileSystemAccessProfile,
} from "@/lib/grpc-client/buildbarn/fsac/fsac";
import type {
  InitialSizeClassCacheClient,
  PreviousExecutionStats,
} from "@/lib/grpc-client/buildbarn/iscc/iscc";
import {
  FilePoolResourceUsage,
  InputRootResourceUsage,
  MonetaryResourceUsage,
  POSIXResourceUsage,
} from "@/lib/grpc-client/buildbarn/resourceusage/resourceusage";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import {
  type BrowserPageParams,
  BrowserPageType,
} from "@/types/BrowserPageType";
import { ProtobufTypeUrls } from "@/types/protobufTypeUrls";
import { getReducedActionDigest_SHA256 } from "@/utils/digestFunctionUtils";
import { fetchCasObject, fetchCasObjectAndParse } from "@/utils/fetchCasObject";
import { protobufToObject } from "@/utils/protobufToObject";
import type { ActionConsoleOutput } from "./types";

export const fetchBrowserActionGrid = async (
  browserPageParams: BrowserPageParams,
  actionCacheClient: ActionCacheClient,
  casByteStreamClient: ByteStreamClient,
  initialSizeClassCacheClient: InitialSizeClassCacheClient,
  fileSystemAccessCacheClient: FileSystemAccessCacheClient,
): Promise<{
  executeResponse: ExecuteResponse | undefined;
  action: Action;
  actionDigest: Digest;
  authenticationMetadata: AuthenticationMetadata | undefined;
  requestMetadata: RequestMetadata | undefined;
  posixResourceUsage: POSIXResourceUsage | undefined;
  filePoolResourceUsage: FilePoolResourceUsage | undefined;
  inputRootResourceUsage: InputRootResourceUsage | undefined;
  monetaryResourceUsage: MonetaryResourceUsage | undefined;
  casCommand: Command | undefined;
  casDirectory: Directory | undefined;
  consoleOutputs: ActionConsoleOutput[];
  previousExecutionStats: PreviousExecutionStats | undefined;
  fileSystemAccessProfile: FileSystemAccessProfile | undefined;
}> => {
  const { actionDigest, executeResponse } = await fetchExecuteResponse(
    browserPageParams,
    casByteStreamClient,
    actionCacheClient,
  );

  const action = await fetchCasObjectAndParse(
    casByteStreamClient,
    browserPageParams.instanceName,
    browserPageParams.digestFunction,
    actionDigest,
    Action,
  );

  const {
    authenticationMetadata,
    requestMetadata,
    posixResourceUsage,
    filePoolResourceUsage,
    inputRootResourceUsage,
    monetaryResourceUsage,
  } = extractMetadataFromExecuteResponse(executeResponse);

  const [
    casCommand,
    casDirectory,
    previousExecutionStats,
    consoleOutputs,
    fileSystemAccessProfile,
  ] = await Promise.all([
    // Fetch Command
    action.commandDigest
      ? fetchCasObjectAndParse(
          casByteStreamClient,
          browserPageParams.instanceName,
          browserPageParams.digestFunction,
          action.commandDigest,
          Command,
        )
      : Promise.resolve(undefined),

    // Fetch Directory
    action.inputRootDigest
      ? fetchCasObjectAndParse(
          casByteStreamClient,
          browserPageParams.instanceName,
          browserPageParams.digestFunction,
          action.inputRootDigest,
          Directory,
        )
      : Promise.resolve(undefined),

    // Fetch Previous Execution Stats
    fetchPreviousExecutionStats(
      action,
      initialSizeClassCacheClient,
      browserPageParams,
    ),

    // Fetch Console Outputs
    getConsoleActionOutputs(
      browserPageParams,
      casByteStreamClient,
      executeResponse,
    ),

    // Fetch File System Access Cache Profile
    fetchFileSystemAccessProfile(
      action,
      fileSystemAccessCacheClient,
      browserPageParams,
    ),
  ]);

  return {
    executeResponse,
    action,
    actionDigest,
    authenticationMetadata,
    requestMetadata,
    posixResourceUsage,
    filePoolResourceUsage,
    inputRootResourceUsage,
    monetaryResourceUsage,
    casCommand,
    casDirectory,
    consoleOutputs,
    previousExecutionStats,
    fileSystemAccessProfile,
  };
};

export const getActionConsoleOutput = async (
  browserPageParams: BrowserPageParams,
  casByteStreamClient: ByteStreamClient,
  digest: Digest | undefined,
  rawBytes: Uint8Array | undefined,
  name: string,
): Promise<ActionConsoleOutput | undefined> => {
  if (rawBytes && rawBytes.length > 0)
    return {
      name,
      digest,
      tooLarge: false,
      notFound: false,
      content: new TextDecoder().decode(rawBytes),
    };

  if (digest === undefined) {
    return undefined;
  }

  const MAX_CONSOLE_OUTPUT_SIZE = 10000;
  if (Number.parseInt(digest.sizeBytes) > MAX_CONSOLE_OUTPUT_SIZE) {
    return {
      name,
      digest,
      tooLarge: true,
      notFound: false,
      content: undefined,
    };
  }

  try {
    const content = await fetchCasObject(
      casByteStreamClient,
      browserPageParams.instanceName,
      browserPageParams.digestFunction,
      digest,
    );
    return {
      name,
      digest,
      tooLarge: false,
      notFound: false,
      content: new TextDecoder().decode(content),
    };
  } catch (e) {
    return {
      name,
      digest,
      tooLarge: false,
      notFound: true,
      content: undefined,
    };
  }
};

async function fetchPreviousExecutionStats(
  action: Action,
  initialSizeClassCacheClient: InitialSizeClassCacheClient,
  browserPageParams: BrowserPageParams,
): Promise<PreviousExecutionStats | undefined> {
  if (!action.commandDigest || !action.platform) {
    return undefined;
  }

  try {
    return await initialSizeClassCacheClient.getPreviousExecutionStats({
      digestFunction: browserPageParams.digestFunction,
      instanceName: browserPageParams.instanceName,
      reducedActionDigest: getReducedActionDigest_SHA256(
        action.commandDigest,
        action.platform,
      ),
    });
  } catch (error) {
    console.log("No previous execution stats found");
  }
}

async function getConsoleActionOutputs(
  browserPageParams: BrowserPageParams,
  casByteStreamClient: ByteStreamClient,
  executeResponse: ExecuteResponse | undefined,
): Promise<ActionConsoleOutput[]> {
  const consoleOutputs: ActionConsoleOutput[] = [];

  const stdoutOutput = await getActionConsoleOutput(
    browserPageParams,
    casByteStreamClient,
    executeResponse?.result?.stdoutDigest,
    executeResponse?.result?.stdoutRaw,
    "Standard output",
  );
  if (stdoutOutput) {
    consoleOutputs.push(stdoutOutput);
  }

  const stderrOutput = await getActionConsoleOutput(
    browserPageParams,
    casByteStreamClient,
    executeResponse?.result?.stderrDigest,
    executeResponse?.result?.stderrRaw,
    "Standard error",
  );
  if (stderrOutput) {
    consoleOutputs.push(stderrOutput);
  }

  return consoleOutputs;
}

function extractMetadataFromExecuteResponse(
  executeResponse: ExecuteResponse | undefined,
): {
  authenticationMetadata: AuthenticationMetadata | undefined;
  requestMetadata: RequestMetadata | undefined;
  posixResourceUsage: POSIXResourceUsage | undefined;
  filePoolResourceUsage: FilePoolResourceUsage | undefined;
  inputRootResourceUsage: InputRootResourceUsage | undefined;
  monetaryResourceUsage: MonetaryResourceUsage | undefined;
} {
  let authenticationMetadata: AuthenticationMetadata | undefined = undefined;
  let requestMetadata: RequestMetadata | undefined = undefined;
  let posixResourceUsage: POSIXResourceUsage | undefined = undefined;
  let filePoolResourceUsage: FilePoolResourceUsage | undefined = undefined;
  let inputRootResourceUsage: InputRootResourceUsage | undefined = undefined;
  let monetaryResourceUsage: MonetaryResourceUsage | undefined = undefined;

  if (!executeResponse?.result?.executionMetadata?.auxiliaryMetadata) {
    return {
      authenticationMetadata,
      requestMetadata,
      posixResourceUsage,
      filePoolResourceUsage,
      inputRootResourceUsage,
      monetaryResourceUsage,
    };
  }

  for (const metadata of executeResponse.result.executionMetadata
    .auxiliaryMetadata) {
    switch (metadata.typeUrl) {
      case ProtobufTypeUrls.AUTHENTICATION_METADATA:
        authenticationMetadata = protobufToObject(
          AuthenticationMetadata,
          metadata.value,
          false,
        );
        break;
      case ProtobufTypeUrls.REQUEST_METADATA:
        requestMetadata = protobufToObject(
          RequestMetadata,
          metadata.value,
          false,
        );
        break;
      case ProtobufTypeUrls.POSIX_RESOURCE_USAGE:
        posixResourceUsage = protobufToObject(
          POSIXResourceUsage,
          metadata.value,
          true,
        );
        break;
      case ProtobufTypeUrls.FILE_POOL_RESOURCE_USAGE:
        filePoolResourceUsage = protobufToObject(
          FilePoolResourceUsage,
          metadata.value,
          true,
        );
        break;
      case ProtobufTypeUrls.INPUT_ROOT_RESOURCE_USAGE:
        inputRootResourceUsage = protobufToObject(
          InputRootResourceUsage,
          metadata.value,
          false,
        );
        break;
      case ProtobufTypeUrls.MONETARY_RESOURCE_USAGE:
        monetaryResourceUsage = protobufToObject(
          MonetaryResourceUsage,
          metadata.value,
          false,
        );
        break;
      default:
        console.error(`Unknown metadata type: ${metadata.typeUrl}`);
        break;
    }
  }

  return {
    authenticationMetadata,
    requestMetadata,
    posixResourceUsage,
    filePoolResourceUsage,
    inputRootResourceUsage,
    monetaryResourceUsage,
  };
}

async function fetchExecuteResponse(
  browserPageParams: BrowserPageParams,
  casByteStreamClient: ByteStreamClient,
  actionCacheClient: ActionCacheClient,
): Promise<{
  actionDigest: Digest;
  executeResponse: ExecuteResponse | undefined;
}> {
  if (
    browserPageParams.browserPageType ===
    BrowserPageType.HistoricalExecuteResponse
  ) {
    const historicalExecuteresponse = await fetchCasObjectAndParse(
      casByteStreamClient,
      browserPageParams.instanceName,
      browserPageParams.digestFunction,
      browserPageParams.digest,
      HistoricalExecuteResponse,
    );

    if (!historicalExecuteresponse.executeResponse?.result) {
      throw new Error(
        "HistoricalExecuteResponse does not contain ExecuteResponse",
      );
    }
    if (!historicalExecuteresponse.actionDigest) {
      throw new Error(
        "HistoricalExecuteResponse does not contain ActionDigest",
      );
    }

    return {
      actionDigest: historicalExecuteresponse.actionDigest,
      executeResponse: historicalExecuteresponse.executeResponse,
    };
  }

  try {
    const actionResult = await actionCacheClient.getActionResult({
      instanceName: browserPageParams.instanceName,
      digestFunction: browserPageParams.digestFunction,
      actionDigest: browserPageParams.digest,
      inlineStdout: true,
      inlineStderr: true,
    });
    return {
      actionDigest: browserPageParams.digest,
      executeResponse: ExecuteResponse.fromPartial({
        result: actionResult,
      }),
    };
  } catch (error) {
    console.log("No execute response was found");
  }

  return { actionDigest: browserPageParams.digest, executeResponse: undefined };
}

async function fetchFileSystemAccessProfile(
  action: Action,
  fileSystemAccessCacheClient: FileSystemAccessCacheClient,
  browserPageParams: BrowserPageParams,
): Promise<FileSystemAccessProfile | undefined> {
  if (!action.commandDigest || !action.platform) {
    return undefined;
  }
  try {
    return await fileSystemAccessCacheClient.getFileSystemAccessProfile({
      digestFunction: browserPageParams.digestFunction,
      instanceName: browserPageParams.instanceName,
      reducedActionDigest: getReducedActionDigest_SHA256(
        action.commandDigest,
        action.platform,
      ),
    });
  } catch (error) {
    console.log("No file system access cache profile was found");
  }
}
