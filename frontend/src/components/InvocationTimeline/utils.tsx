import { INVOCATION_RESULT_TAGS } from "../InvocationResultTag";
import { getInvocationResultTagEnum } from "../InvocationResultTag/enum";

export const getInvocationResultTagColor = (
  exitCodeName: string | undefined,
  bepCompleted: boolean,
): string => {
  return INVOCATION_RESULT_TAGS[
    getInvocationResultTagEnum(exitCodeName, bepCompleted)
  ].color;
};
