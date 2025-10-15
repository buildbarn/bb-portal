import type { FilterValue } from "antd/es/table/interface";
import type { BazelInvocationWhereInput } from "@/graphql/__generated__/graphql";
import { InvocationExitCodes, InvocationResultTagEnum } from "./enum";

export const invocationResultTagFilters = [
  { text: "Succeeded", value: InvocationResultTagEnum.SUCCESS },
  { text: "Unstable", value: InvocationResultTagEnum.UNSTABLE },
  { text: "Parsing Failed", value: InvocationResultTagEnum.PARSING_FAILURE },
  { text: "Build Failed", value: InvocationResultTagEnum.BUILD_FAILURE },
  { text: "Tests Failed", value: InvocationResultTagEnum.TESTS_FAILED },
  { text: "Remote error", value: InvocationResultTagEnum.REMOTE_ERROR },
  { text: "Not Built", value: InvocationResultTagEnum.NOT_BUILT },
  { text: "Aborted", value: InvocationResultTagEnum.ABORTED },
  { text: "Interrupted", value: InvocationResultTagEnum.INTERRUPTED },
  { text: "Status Unknown", value: InvocationResultTagEnum.UNKNOWN },
  { text: "In Progress", value: InvocationResultTagEnum.IN_PROGRESS },
  {
    text: "BEP Upload Aborted",
    value: InvocationResultTagEnum.BEP_UPLOAD_ABORTED,
  },
];
