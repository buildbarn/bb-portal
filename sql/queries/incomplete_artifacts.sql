-- name: CreateIncompleteArtifactGraphs :exec
-- Bulk-insert serialized BuildEvents (NamedSetOfFiles / TargetCompleted)
-- into the artifact-graph staging table as they stream in.
INSERT INTO incomplete_artifact_graphs (
    bazel_invocation_id,
    seq_id,
    event
)
SELECT
    sqlc.arg(bazel_invocation_id),
    seq_id,
    event
FROM (
    SELECT
        unnest(sqlc.arg(seq_ids)::int[]) AS seq_id,
        unnest(sqlc.arg(events)::bytea[]) AS event
) AS input;

-- name: GetIncompleteArtifactGraphEvents :many
-- Read all staged artifact-graph events for an invocation, in arrival
-- order. Used both for compaction and for serving the graph in its
-- partial state before compaction has run.
SELECT event
FROM incomplete_artifact_graphs
WHERE bazel_invocation_id = sqlc.arg(bazel_invocation_id)
ORDER BY seq_id;

-- name: DeleteIncompleteArtifactGraphsFromPages :execrows
-- Delete staged events for invocations whose graph has already been
-- compacted into invocation_artifact_graphs. Paged via ctid so cleanup is
-- spread across cleanup ticks via nextSlice(), mirroring
-- DeleteIncompleteLogsFromPages.
DELETE FROM incomplete_artifact_graphs g
USING bazel_invocations AS i
WHERE
    g.bazel_invocation_id = i.id
    AND g.ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
    AND g.ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
    AND i.bep_completed = true
    AND EXISTS (
        SELECT 1 FROM invocation_artifact_graphs a
        WHERE a.bazel_invocation_id = i.id
    );
