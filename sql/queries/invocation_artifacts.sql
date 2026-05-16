-- name: InsertInvocationArtifactGraph :exec
-- Insert (or replace) the compressed artifact graph for an invocation.
-- Idempotent via the unique edge on bazel_invocation_artifact_graph.
INSERT INTO invocation_artifact_graphs (
    payload,
    uncompressed_size,
    bazel_invocation_artifact_graph
) VALUES (
    sqlc.arg(payload)::bytea,
    sqlc.arg(uncompressed_size)::bigint,
    sqlc.arg(bazel_invocation_id)::bigint
)
ON CONFLICT (bazel_invocation_artifact_graph) DO UPDATE SET
    payload = EXCLUDED.payload,
    uncompressed_size = EXCLUDED.uncompressed_size;

-- name: GetInvocationArtifactGraph :one
SELECT payload, uncompressed_size
FROM invocation_artifact_graphs
WHERE bazel_invocation_artifact_graph = sqlc.arg(bazel_invocation_id);

-- name: DeleteOldInvocationArtifactGraphsFromPages :execrows
-- Delete artifact-graph rows whose owning invocation completed more than
-- artifact_retention ago. Paged via ctid on the artifact graphs table —
-- same pattern as DeleteOldInvocationsFromPages so cleanup is evenly
-- spread across cleanup ticks via nextSlice().
DELETE FROM invocation_artifact_graphs
WHERE
    ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
    AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
    AND bazel_invocation_artifact_graph IN (
      SELECT id FROM bazel_invocations
      WHERE bep_completed = true
        AND ended_at IS NOT NULL
        AND ended_at < sqlc.arg(cutoff_time)::timestamptz
    );
