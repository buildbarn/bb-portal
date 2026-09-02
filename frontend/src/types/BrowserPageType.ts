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
