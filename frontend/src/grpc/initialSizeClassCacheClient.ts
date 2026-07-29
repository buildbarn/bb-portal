import { createChannel, createClient } from "nice-grpc-web";
import {
  type InitialSizeClassCacheClient,
  InitialSizeClassCacheDefinition,
} from "@/lib/grpc-client/buildbarn/iscc/iscc";

const grpcChannel = createChannel("/api/v1/grpcweb");

export const initialSizeClassCacheClient: InitialSizeClassCacheClient =
  createClient(InitialSizeClassCacheDefinition, grpcChannel);
