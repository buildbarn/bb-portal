-- name: CreateTestSummariesBulk :execrows
INSERT INTO test_summaries (
    invocation_target_test_summary
)
SELECT
    it.id
FROM (
    -- STAGE 1: Resolve Target IDs immediately using the Index
    SELECT
        input.config_id,
        t.id AS target_id
    FROM (
        SELECT
            sqlc.arg(instance_name_id)::bigint as instance_name_id,
            unnest(sqlc.arg(labels)::text[]) as label,
            unnest(sqlc.arg(config_ids)::text[]) as config_id
    ) AS input
    JOIN targets t
        ON t.instance_name_targets = input.instance_name_id
        AND t.label = input.label
        AND t.aspect = ''
    -- We use OFFSET 0 here to force this block to execute first, which is alot faster
    OFFSET 0
) AS resolved_targets
-- STAGE 2: Join the rest using the specific Target IDs we found
JOIN configurations c
    ON c.bazel_invocation_id = sqlc.arg(bazel_invocation_id)
    AND c.configuration_id = resolved_targets.config_id
JOIN invocation_targets it
    ON it.target_invocation_targets = resolved_targets.target_id
    AND it.invocation_target_configuration = c.id
    AND it.bazel_invocation_invocation_targets = sqlc.arg(bazel_invocation_id);

-- name: UpdateTestSummariesBulk :execrows
UPDATE test_summaries ts
SET
    overall_status = source.overall_status,
    total_run_count = source.total_run_count,
    run_count = source.run_count,
    attempt_count = source.attempt_count,
    shard_count = source.shard_count,
    total_num_cached = source.total_num_cached,
    first_start_time = source.first_start_time,
    last_stop_time = source.last_stop_time,
    total_run_duration_in_ms = source.total_run_duration_in_ms
FROM (
    SELECT
        it.id AS invocation_target_id,
        resolved_targets.overall_status,
        resolved_targets.total_run_count,
        resolved_targets.run_count,
        resolved_targets.attempt_count,
        resolved_targets.shard_count,
        resolved_targets.total_num_cached,
        resolved_targets.first_start_time,
        resolved_targets.last_stop_time,
        resolved_targets.total_run_duration_in_ms
    FROM (
        -- STAGE 1: Resolve Target IDs immediately using the Index
        SELECT
            input.config_id,
            input.overall_status,
            input.total_run_count,
            input.run_count,
            input.attempt_count,
            input.shard_count,
            input.total_num_cached,
            input.first_start_time,
            input.last_stop_time,
            input.total_run_duration_in_ms,
            t.id AS target_id
        FROM (
            SELECT
                sqlc.arg(instance_name_id)::bigint as instance_name_id,
                unnest(sqlc.arg(labels)::text[]) as label,
                unnest(sqlc.arg(config_ids)::text[]) as config_id,
                unnest(sqlc.arg(overall_statuses)::text[]) as overall_status,
                unnest(sqlc.arg(total_run_counts)::int[]) as total_run_count,
                unnest(sqlc.arg(run_counts)::int[]) as run_count,
                unnest(sqlc.arg(attempt_counts)::int[]) as attempt_count,
                unnest(sqlc.arg(shard_counts)::int[]) as shard_count,
                unnest(sqlc.arg(total_num_cacheds)::int[]) as total_num_cached,
                unnest(sqlc.arg(start_times)::timestamptz[]) as first_start_time,
                unnest(sqlc.arg(stop_times)::timestamptz[]) as last_stop_time,
                unnest(sqlc.arg(durations)::bigint[]) as total_run_duration_in_ms
        ) AS input
        JOIN targets t
            ON t.instance_name_targets = input.instance_name_id
            AND t.label = input.label
            AND t.aspect = ''
        -- Optimization: Force this block to execute first
        OFFSET 0
    ) AS resolved_targets
    -- STAGE 2: Join the rest using the specific Target IDs we found
    JOIN configurations c
        ON c.bazel_invocation_id = sqlc.arg(bazel_invocation_id)
        AND c.configuration_id = resolved_targets.config_id
    JOIN invocation_targets it
        ON it.target_invocation_targets = resolved_targets.target_id
        AND it.invocation_target_configuration = c.id
        AND it.bazel_invocation_invocation_targets = sqlc.arg(bazel_invocation_id)
) AS source
WHERE 
    ts.invocation_target_test_summary = source.invocation_target_id;
