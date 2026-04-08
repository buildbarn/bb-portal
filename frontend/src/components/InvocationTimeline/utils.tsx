import { INVOCATION_RESULT_TAGS } from "../InvocationResultTag";
import { getInvocationResultTagEnum } from "../InvocationResultTag/enum";

export const getInvocationResultTagColor = (
  exitCodeName: string | undefined,
  timeSinceLastConnectionMillis: number | undefined,
): string => {
  return INVOCATION_RESULT_TAGS[
    getInvocationResultTagEnum(exitCodeName, timeSinceLastConnectionMillis)
  ].color;
};
