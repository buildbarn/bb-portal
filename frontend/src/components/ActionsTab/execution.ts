export type ActionExecutionKind = "Remote" | "Local" | "Internal" | "Unknown";

export const getActionExecutionKind = (
  runner: string | null | undefined,
): ActionExecutionKind => {
  switch (runner) {
    case "remote":
    case "remote cache hit":
      return "Remote";
    case "internal":
      return "Internal";
    case null:
    case undefined:
    case "":
      return "Unknown";
    default:
      return "Local";
  }
};
