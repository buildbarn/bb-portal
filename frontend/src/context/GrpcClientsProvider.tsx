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
  const buildQueueStateClient: BuildQueueStateClient = createClient(
    BuildQueueStateDefinition,
    createChannel(env("NEXT_PUBLIC_BB_BUILDQUEUESTATE_GRPC_BACKEND_URL") || ""),
  );

  const actionCacheClient: ActionCacheClient = createClient(
    ActionCacheDefinition,
    createChannel(env("NEXT_PUBLIC_BB_ACTIONCACHE_GRPC_BACKEND_URL") || ""),
  );

  const casByteStreamClient: ByteStreamClient = createClient(
    ByteStreamDefinition,
    createChannel(env("NEXT_PUBLIC_BB_CAS_GRPC_BACKEND_URL") || ""),
  );

  const initialSizeClassCacheClient: InitialSizeClassCacheClient = createClient(
    InitialSizeClassCacheDefinition,
    createChannel(env("NEXT_PUBLIC_BB_ISCC_GRPC_BACKEND_URL") || ""),
  );

  const fileSystemAccessCacheClient: FileSystemAccessCacheClient = createClient(
    FileSystemAccessCacheDefinition,
    createChannel(env("NEXT_PUBLIC_BB_FSAC_GRPC_BACKEND_URL") || ""),
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
