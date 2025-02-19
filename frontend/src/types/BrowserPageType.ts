import type {
  Digest,
  DigestFunction_Value,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

export enum BrowserPageType {
  Action = "action",
  Command = "command",
  Directory = "directory",
  File = "file",
  HistoricalExecuteResponse = "historical_execute_response",
  PreviousExecutionStats = "previous_execution_stats",
  Tree = "tree",
}

export const getBrowserPageTypeFromString = (
  value: string,
): BrowserPageType | undefined => {
  if (Object.values(BrowserPageType).includes(value as BrowserPageType)) {
    return value as BrowserPageType;
  }
  return undefined;
};

export interface BrowserPageParams {
  instanceName: string;
  digestFunction: DigestFunction_Value;
  browserPageType: BrowserPageType;
  digest: Digest;
  otherParams: Array<string>;
}
