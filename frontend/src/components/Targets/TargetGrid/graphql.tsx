import { gql } from "@/graphql/__generated__";

export const GET_TARGETS = gql(/* GraphQl */`
query GetUniqueTargetLabels($label: String) {
  getUniqueTargetLabels(param: $label)
}
`);