-- name: CreateIncompleteBuildLogs :exec
INSERT INTO incomplete_build_logs (
    bazel_invocation_id,
    snippet_id,
    log_snippet
)
SELECT 
    sqlc.arg(bazel_invocation_id),
    snippet_id,
    log_snippet
FROM (
    SELECT 
        unnest(sqlc.arg(snippet_ids)::int[]) AS snippet_id,
        unnest(sqlc.arg(log_snippets)::bytea[]) AS log_snippet
) AS input;
