import type { DocumentTypeDecoration } from "@graphql-typed-document-node/core";
import { type FragmentType, getFragmentData } from "@/graphql/__generated__";

interface GraphqlEdgeList<NodeType> {
  edges?: Array<{ node?: NodeType | null } | null | undefined> | null;
}

export function parseGraphqlEdgeList<NodeType>(
  data: GraphqlEdgeList<NodeType> | undefined | null,
): NodeType[] {
  if (!data || !data.edges) {
    return [];
  }
  return data.edges.reduce<NodeType[]>((acc, edge) => {
    if (!edge?.node) {
      return acc;
    }
    acc.push(edge.node);
    return acc;
  }, []);
}

export function parseGraphqlEdgeListWithFragment<TType>(
  fragment_definition: DocumentTypeDecoration<TType, any>,
  data:
    | GraphqlEdgeList<
        FragmentType<DocumentTypeDecoration<TType, any>> | null | undefined
      >
    | undefined
    | null,
): TType[] {
  if (!data || !data.edges) {
    return [];
  }
  return data.edges.reduce<TType[]>((acc, edge) => {
    if (!edge?.node) {
      return acc;
    }
    const fragmentData = getFragmentData(fragment_definition, edge.node);
    if (fragmentData) {
      acc.push(fragmentData);
    }
    return acc;
  }, []);
}
