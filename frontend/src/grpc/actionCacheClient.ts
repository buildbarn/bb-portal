import { createChannel, createClient } from "nice-grpc-web";
import {
  type ActionCacheClient,
  ActionCacheDefinition,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

const grpcChannel = createChannel("/api/v1/grpcweb");
export const actionCacheClient: ActionCacheClient = createClient(
  ActionCacheDefinition,
  grpcChannel,
);
