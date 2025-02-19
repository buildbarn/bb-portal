import type { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

export type ActionConsoleOutput = {
  name: string;
  digest: Digest | undefined;
  tooLarge: boolean;
  notFound: boolean;
  content: string | undefined;
};
