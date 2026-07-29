import type { Digest } from "../lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

// This interface was originally generated from query.proto in bb-browser,
// but has been implemented here manually since the bb-browser dependency
// has been removed.
export interface FileSystemAccessProfileReference {
  digest: Digest | undefined;
  pathHashesBaseHash: string;
}
