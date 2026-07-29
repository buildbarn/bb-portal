import { createChannel, createClient } from "nice-grpc-web";
import {
  type ByteStreamClient,
  ByteStreamDefinition,
} from "@/lib/grpc-client/google/bytestream/bytestream";

const grpcChannel = createChannel("/api/v1/grpcweb");
export const casByteStreamClient: ByteStreamClient = createClient(
  ByteStreamDefinition,
  grpcChannel,
);
