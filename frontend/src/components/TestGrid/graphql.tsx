import { gql } from "@/graphql/__generated__";

export const GET_TEST_LABELS = gql(/* GraphQl */`
query GetUniqueTestLabels($label: String) {
  getUniqueTestLabels(param: $label)
}
`);