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

-- name: DeleteIncompleteLogsFromPages :execrows
DELETE FROM incomplete_build_logs
WHERE ctid IN (
    SELECT l.ctid
    FROM incomplete_build_logs l
    JOIN bazel_invocations AS i ON i.id = l.bazel_invocation_id
    WHERE
        l.bazel_invocation_id = i.id
        AND l.ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
        AND l.ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
        AND i.bep_completed = true
        AND EXISTS (
            SELECT 1 FROM build_log_chunks c
            WHERE c.bazel_invocation_build_log_chunks = i.id
        )
    FOR UPDATE SKIP LOCKED
    LIMIT sqlc.arg(batch_limit)::bigint
);
