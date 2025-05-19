import {
  type ActionCacheClient,
  ActionCacheDefinition,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import {
  type BuildQueueStateClient,
  BuildQueueStateDefinition,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import {
  type FileSystemAccessCacheClient,
  FileSystemAccessCacheDefinition,
} from "@/lib/grpc-client/buildbarn/fsac/fsac";
import {
  type InitialSizeClassCacheClient,
  InitialSizeClassCacheDefinition,
} from "@/lib/grpc-client/buildbarn/iscc/iscc";
import {
  type ByteStreamClient,
  ByteStreamDefinition,
} from "@/lib/grpc-client/google/bytestream/bytestream";
import { env } from "next-runtime-env";
import { createChannel, createClient } from "nice-grpc-web";
import type { ReactNode } from "react";
import { GrpcClientsContext } from "./GrpcClientsContext";

export interface GrpcClientsProviderProps {
  children: ReactNode;
}

const GrpcClientsProvider = ({ children }: GrpcClientsProviderProps) => {
  const grpcChannel = createChannel(
    `${env("NEXT_PUBLIC_BES_BACKEND_URL") || ""}/api/v1/grpcweb`,
  );

  const buildQueueStateClient: BuildQueueStateClient = createClient(
    BuildQueueStateDefinition,
    grpcChannel,
  );

  const actionCacheClient: ActionCacheClient = createClient(
    ActionCacheDefinition,
    grpcChannel,
  );

  const casByteStreamClient: ByteStreamClient = createClient(
    ByteStreamDefinition,
    grpcChannel,
  );

  const initialSizeClassCacheClient: InitialSizeClassCacheClient = createClient(
    InitialSizeClassCacheDefinition,
    grpcChannel,
  );

  const fileSystemAccessCacheClient: FileSystemAccessCacheClient = createClient(
    FileSystemAccessCacheDefinition,
    grpcChannel,
  );

  return (
    <GrpcClientsContext.Provider
      value={{
        buildQueueStateClient,
        actionCacheClient,
        casByteStreamClient,
        initialSizeClassCacheClient,
        fileSystemAccessCacheClient,
      }}
    >
      {children}
    </GrpcClientsContext.Provider>
  );
};

export default GrpcClientsProvider;
