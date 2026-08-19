-- name: GetActionTimingMetrics :one
WITH action_durations AS (
    SELECT
        cache_hit,
        CASE
            WHEN start_time IS NOT NULL
                AND end_time IS NOT NULL
                AND end_time >= start_time
            THEN EXTRACT(EPOCH FROM (end_time - start_time)) * 1000
        END AS duration_in_ms
    FROM actions
    WHERE bazel_invocation_id = $1
)
SELECT
    COALESCE(FLOOR(SUM(duration_in_ms)), 0)::bigint AS total_expected_time_in_ms,
    COALESCE(
        FLOOR(SUM(duration_in_ms) FILTER (WHERE cache_hit IS TRUE)),
        0
    )::bigint AS time_saved_by_cache_hits_in_ms,
    COUNT(*)::bigint AS total_actions,
    COUNT(duration_in_ms)::bigint AS timed_actions,
    COUNT(*) FILTER (WHERE cache_hit IS TRUE)::bigint AS cache_hit_actions,
    COUNT(duration_in_ms) FILTER (WHERE cache_hit IS TRUE)::bigint AS timed_cache_hit_actions
FROM action_durations;
