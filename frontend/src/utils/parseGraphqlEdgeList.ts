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
    if (edge?.node) {
      acc.push(edge.node);
    }
    return acc;
  }, []);
}
