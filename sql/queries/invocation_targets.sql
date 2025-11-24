-- name: CreateInvocationTargetsBulk :exec
INSERT INTO invocation_targets (
    bazel_invocation_invocation_targets,
    target_invocation_targets,
    success,
    tags,
    start_time_in_ms,
    end_time_in_ms,
    duration_in_ms,
    failure_message,
    abort_reason
)
SELECT
    sqlc.arg(bazel_invocation_id),
    target_id,
    success,
    NULLIF(tags, '')::jsonb,
    NULLIF(start_time, 0),
    NULLIF(end_time, 0),
    NULLIF(duration, 0),
    NULLIF(failure_message, ''),
    abort_reason
FROM (
    SELECT 
        unnest(sqlc.arg(target_ids)::bigint[]) AS target_id,
        unnest(sqlc.arg(successes)::boolean[]) AS success,
        unnest(sqlc.arg(tags_list)::text[]) AS tags,
        unnest(sqlc.arg(start_times)::bigint[]) AS start_time,
        unnest(sqlc.arg(end_times)::bigint[]) AS end_time,
        unnest(sqlc.arg(durations)::bigint[]) AS duration,
        unnest(sqlc.arg(failure_messages)::text[]) AS failure_message,
        unnest(sqlc.arg(abort_reasons)::text[]) AS abort_reason
) AS input;
