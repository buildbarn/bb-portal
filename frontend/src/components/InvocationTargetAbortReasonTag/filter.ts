import { InvocationTargetAbortReason } from "@/graphql/__generated__/graphql";

export const getInvocationTargetAbortReasonFilterOptions = () => {
  return Object.entries(InvocationTargetAbortReason).map(([key, value]) => ({
    text: key,
    value: value as InvocationTargetAbortReason,
  }));
};
