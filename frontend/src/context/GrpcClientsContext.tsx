import type {
  ActionCacheClient,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BuildQueueStateClient } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { FileSystemAccessCacheClient } from "@/lib/grpc-client/buildbarn/fsac/fsac";
import type { InitialSizeClassCacheClient } from "@/lib/grpc-client/buildbarn/iscc/iscc";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import { createContext, useContext } from "react";

export type CasObjectFetchFunction = <T>(
  objectType: {
    decode: (input: Uint8Array) => T;
    toJSON: (input: T) => unknown;
  },
  instanceName: string | undefined,
  digestFunction: DigestFunction_Value,
  digest: string,
  sizeBytes: string,
) => Promise<T>;

interface GrpcClientsContextState {
  buildQueueStateClient: BuildQueueStateClient;
  actionCacheClient: ActionCacheClient;
  casByteStreamClient: ByteStreamClient;
  initialSizeClassCacheClient: InitialSizeClassCacheClient;
  fileSystemAccessCacheClient: FileSystemAccessCacheClient;
}

// biome-ignore lint/style/noNonNullAssertion: We want to throw an error if the context is used without provider, instead of failing silently.
export const GrpcClientsContext = createContext<GrpcClientsContextState>(null!);

export const useGrpcClients = () => useContext(GrpcClientsContext);
