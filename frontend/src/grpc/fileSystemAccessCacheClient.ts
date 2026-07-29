import { createChannel, createClient } from "nice-grpc-web";
import {
  type FileSystemAccessCacheClient,
  FileSystemAccessCacheDefinition,
} from "@/lib/grpc-client/buildbarn/fsac/fsac";

const grpcChannel = createChannel("/api/v1/grpcweb");

export const fileSystemAccessCacheClient: FileSystemAccessCacheClient =
  createClient(FileSystemAccessCacheDefinition, grpcChannel);
