// The artifact graph is decoded server-side and delivered as structured
// GraphQL data (see internal/graphql/artifact_graph.go). The client only
// indexes the named sets by id and walks the graph to enumerate the files
// reachable from a target's output groups — no decompression or wire
// parsing happens here.

export interface ArtifactFileNode {
  name: string;
  uri?: string | null;
  digest?: string | null;
  sizeBytes?: number | null;
  downloadUrl?: string | null;
}

export interface NamedSetNode {
  id: string;
  files: ArtifactFileNode[];
  childSetIds: string[];
}

export interface OutputGroupNode {
  name: string;
  incomplete: boolean;
  rootSetIds: string[];
}

export interface TargetNode {
  label: string;
  aspect?: string | null;
  outputGroups: OutputGroupNode[];
}

export interface ArtifactGraphData {
  namedSets: NamedSetNode[];
  targets: TargetNode[];
}

// buildSetIndex indexes named sets by their id for graph traversal.
export function buildSetIndex(
  namedSets: readonly NamedSetNode[],
): Map<string, NamedSetNode> {
  const index = new Map<string, NamedSetNode>();
  for (const set of namedSets) {
    index.set(set.id, set);
  }
  return index;
}

// walkSetFiles walks the file graph rooted at the given set ids and yields
// every reachable file, terminating on cycles.
export function* walkSetFiles(
  sets: ReadonlyMap<string, NamedSetNode>,
  rootSetIds: readonly string[],
): Generator<ArtifactFileNode> {
  const visited = new Set<string>();
  const queue: string[] = [...rootSetIds];
  while (queue.length > 0) {
    const id = queue.shift() as string;
    if (visited.has(id)) continue;
    visited.add(id);
    const set = sets.get(id);
    if (!set) continue;
    for (const file of set.files) {
      yield file;
    }
    for (const childId of set.childSetIds) {
      queue.push(childId);
    }
  }
}

// countSetFiles returns the number of distinct files reachable from the
// given root set ids.
export function countSetFiles(
  sets: ReadonlyMap<string, NamedSetNode>,
  rootSetIds: readonly string[],
): number {
  let n = 0;
  for (const _ of walkSetFiles(sets, rootSetIds)) {
    n++;
  }
  return n;
}
