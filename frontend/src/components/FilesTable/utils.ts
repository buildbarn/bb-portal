import type {
  ActionResult,
  Command,
  DigestFunction_Value,
  LogFile,
  OutputDirectory,
  OutputFile,
  OutputSymlink,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { generateDirectoryUrl, generateFileUrl } from "@/utils/urlGenerator";
import type { FilesTableEntry } from "./Columns";

export function filesTableEntryFromOutputDirectory(
  outputDirectory: OutputDirectory,
  instanceName: string,
  digestFunction: DigestFunction_Value,
): FilesTableEntry {
  const digest = outputDirectory.rootDirectoryDigest
    ? outputDirectory.rootDirectoryDigest
    : outputDirectory.treeDigest;

  return {
    mode: "drwxr-xr-x",
    size: digest?.sizeBytes,
    filename: outputDirectory.path,
    href: digest
      ? generateDirectoryUrl(instanceName, digestFunction, digest)
      : undefined,
  };
}

export function filesTableEntryFromOutputSymlink(
  outputSymlink: OutputSymlink,
): FilesTableEntry {
  return {
    mode: "lrwxrwxrwx",
    size: undefined,
    filename: `${outputSymlink.path} -> ${outputSymlink.target}`,
    href: undefined,
  };
}

export function filesTableEntryFromOutputFile(
  outputFile: OutputFile,
  instanceName: string,
  digestFunction: DigestFunction_Value,
): FilesTableEntry {
  return {
    mode: `-rw${outputFile.isExecutable ? "x" : "-"}r-${
      outputFile.isExecutable ? "x" : "-"
    }r-${outputFile.isExecutable ? "x" : "-"}`,
    size: outputFile.digest?.sizeBytes,
    filename: outputFile.path,
    href: outputFile.digest
      ? generateFileUrl(
          instanceName,
          digestFunction,
          outputFile.digest,
          outputFile.path.split("/").slice(-1)[0],
        )
      : undefined,
  };
}

export function filesTableEntriesFromOutputPath(
  outputPath: string,
): FilesTableEntry {
  return {
    mode: undefined,
    size: undefined,
    filename: outputPath,
    href: undefined,
  };
}

export function filesTableEntriesFromActionResultAndCommand(
  actionResult: ActionResult | undefined,
  command: Command | undefined,
  instanceName: string,
  digestFunction: DigestFunction_Value,
): FilesTableEntry[] {
  const entries: FilesTableEntry[] = [];

  if (actionResult) {
    if (actionResult.outputDirectories) {
      for (const outputDirectory of actionResult.outputDirectories) {
        entries.push(
          filesTableEntryFromOutputDirectory(
            outputDirectory,
            instanceName,
            digestFunction,
          ),
        );
      }
    }

    if (actionResult.outputSymlinks) {
      for (const outputSymlink of actionResult.outputSymlinks) {
        entries.push(filesTableEntryFromOutputSymlink(outputSymlink));
      }
    }
    if (actionResult.outputFiles) {
      for (const outputFile of actionResult.outputFiles) {
        entries.push(
          filesTableEntryFromOutputFile(
            outputFile,
            instanceName,
            digestFunction,
          ),
        );
      }
    }
  }

  if (command) {
    for (const outputPath of command.outputPaths) {
      if (
        !entries.find((filesTableEntry) => {
          return filesTableEntry.filename === outputPath;
        })
      ) {
        entries.push(filesTableEntriesFromOutputPath(outputPath));
      }
    }
  }

  return entries;
}

export function filesTableEntriesFromServerLogs(
  serverLogs: {
    [key: string]: LogFile;
  },
  instanceName: string,
  digestFunction: DigestFunction_Value,
): FilesTableEntry[] {
  const entries: FilesTableEntry[] = [];

  for (const key of Object.keys(serverLogs)) {
    const logFile = serverLogs[key];
    entries.push({
      mode: "-rw-r--r--",
      size: logFile.digest?.sizeBytes,
      filename: key,
      href: logFile.digest
        ? generateFileUrl(instanceName, digestFunction, logFile.digest, key)
        : undefined,
    });
  }

  return entries;
}
