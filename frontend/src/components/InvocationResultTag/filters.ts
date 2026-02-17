import type { FilterValue } from "antd/es/table/interface";
import type { BazelInvocationWhereInput } from "@/graphql/__generated__/graphql";
import {
  INVOCATION_IN_PROGRESS_TIMEOUT,
  InvocationExitCodes,
  InvocationResult,
} from "./enum";

export const invocationResultTagFilters = [
  { text: "Succeeded", value: InvocationResult.SUCCESS },
  { text: "Unstable", value: InvocationResult.UNSTABLE },
  { text: "Parsing Failed", value: InvocationResult.PARSING_FAILURE },
  { text: "Build Failed", value: InvocationResult.BUILD_FAILURE },
  { text: "Tests Failed", value: InvocationResult.TESTS_FAILED },
  { text: "Remote error", value: InvocationResult.REMOTE_ERROR },
  { text: "Not Built", value: InvocationResult.NOT_BUILT },
  { text: "Aborted", value: InvocationResult.ABORTED },
  { text: "Interrupted", value: InvocationResult.INTERRUPTED },
  { text: "Unknown exit code", value: InvocationResult.UNKNOWN_EXIT_CODE },
  { text: "In Progress", value: InvocationResult.IN_PROGRESS },
  { text: "Disconnected", value: InvocationResult.DISCONNECTED },
];

export const applyInvocationResultTagFilter = (
  value: FilterValue,
): BazelInvocationWhereInput[] => {
  const filters: BazelInvocationWhereInput[] = [];
  value.forEach((v) => {
    const connectionCutoff = new Date(
      Date.now() - INVOCATION_IN_PROGRESS_TIMEOUT,
    );
    const tag = v as InvocationResult;
    switch (tag) {
      case InvocationResult.IN_PROGRESS:
        filters.push({
          and: [
            { or: [{ exitCodeNameIsNil: true }, { exitCodeName: "" }] },
            {
              hasConnectionMetadataWith: [
                { connectionLastOpenAtGTE: connectionCutoff },
              ],
            },
          ],
        });
        break;
      case InvocationResult.DISCONNECTED:
        filters.push({
          and: [
            { or: [{ exitCodeNameIsNil: true }, { exitCodeName: "" }] },
            {
              hasConnectionMetadataWith: [
                { connectionLastOpenAtLT: connectionCutoff },
              ],
            },
          ],
        });
        break;
      case InvocationResult.UNKNOWN_EXIT_CODE:
        filters.push({
          and: [
            { exitCodeNameNotIn: Object.values(InvocationExitCodes) },
            { exitCodeNameNotNil: true },
            { exitCodeNameNEQ: "" },
          ],
        });
        break;
      default:
        filters.push({ exitCodeName: tag });
        break;
    }
  });
  return [{ or: filters }];
};
