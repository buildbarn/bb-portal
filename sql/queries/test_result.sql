-- name: CreateTestResultsBulk :many
INSERT INTO test_results (
    test_summary_test_results,
    run,
    shard,
    attempt,
    status,
    status_details,
    cached_locally,
    test_attempt_start,
    test_attempt_duration_in_ms,
    warning,
    strategy,
    cached_remotely,
    exit_code,
    hostname,
    timing_breakdown
)
SELECT
    ts.id,
    resolved_inputs.run,
    resolved_inputs.shard,
    resolved_inputs.attempt,
    resolved_inputs.status,
    resolved_inputs.status_details,
    resolved_inputs.cached_locally,
    resolved_inputs.test_attempt_start,
    resolved_inputs.test_attempt_duration_in_ms,
    resolved_inputs.warning::jsonb,
    resolved_inputs.strategy,
    resolved_inputs.cached_remotely,
    resolved_inputs.exit_code,
    resolved_inputs.hostname,
    resolved_inputs.timing_breakdown::jsonb
FROM (
    -- STAGE 1: Resolve Target IDs immediately using the Index
    -- We join the massive input arrays against the Targets table first.
    SELECT
        input.config_id,
        t.id AS target_id,
        input.run,
        input.shard,
        input.attempt,
        input.status,
        input.status_details,
        input.cached_locally,
        input.test_attempt_start,
        input.test_attempt_duration_in_ms,
        input.warning,
        input.strategy,
        input.cached_remotely,
        input.exit_code,
        input.hostname,
        input.timing_breakdown,
        input.input_order
    -- The caller matches the returned IDs to its inputs by position, so
    -- number the inputs here, before any join: a join may emit rows in any
    -- order (e.g. a hash join emits them in targets order).
    FROM ROWS FROM (
        unnest(sqlc.arg(labels)::text[]),
        unnest(sqlc.arg(config_ids)::text[]),
        unnest(sqlc.arg(runs)::integer[]),
        unnest(sqlc.arg(shards)::integer[]),
        unnest(sqlc.arg(attempts)::integer[]),
        unnest(sqlc.arg(statuses)::text[]),
        unnest(sqlc.arg(status_detailss)::text[]),
        unnest(sqlc.arg(cached_locallys)::boolean[]),
        unnest(sqlc.arg(test_attempt_starts)::timestamptz[]),
        unnest(sqlc.arg(test_attempt_durations)::bigint[]),
        unnest(sqlc.arg(warnings)::text[]),
        unnest(sqlc.arg(strategies)::text[]),
        unnest(sqlc.arg(cached_remotelys)::boolean[]),
        unnest(sqlc.arg(exit_codes)::integer[]),
        unnest(sqlc.arg(hostnames)::text[]),
        unnest(sqlc.arg(timing_breakdowns)::text[])
    ) WITH ORDINALITY AS input(
        label,
        config_id,
        run,
        shard,
        attempt,
        status,
        status_details,
        cached_locally,
        test_attempt_start,
        test_attempt_duration_in_ms,
        warning,
        strategy,
        cached_remotely,
        exit_code,
        hostname,
        timing_breakdown,
        input_order
    )
    JOIN targets t
        ON t.instance_name_targets = sqlc.arg(instance_name_id)::bigint
        AND t.label = input.label
        AND t.aspect = ''
    -- We use OFFSET 0 here to force this block to execute first
    OFFSET 0
) AS resolved_inputs
-- STAGE 2: Join the rest using the specific Target IDs we found
JOIN configurations c
    ON c.bazel_invocation_id = sqlc.arg(bazel_invocation_id)
    AND c.configuration_id = resolved_inputs.config_id
JOIN invocation_targets it
    ON it.target_invocation_targets = resolved_inputs.target_id
    AND it.invocation_target_configuration = c.id
    AND it.bazel_invocation_invocation_targets = sqlc.arg(bazel_invocation_id)
JOIN test_summaries ts
    ON ts.invocation_target_test_summary = it.id
-- Order the results by the input order
ORDER BY resolved_inputs.input_order
RETURNING id;

-- name: CreateTestActionOutputAssociation :execrows
INSERT INTO test_action_outputs (
    test_result_id,
    file_id
)
SELECT
    input.test_result_id,
    f.id
FROM (
    SELECT
        unnest(sqlc.arg(test_result_id)::bigint[]) AS test_result_id,
        unnest(sqlc.arg(file_paths)::text[]) AS file_path,
        unnest(sqlc.arg(rev2_instance_names)::text[]) AS rev2_instance_name,
        unnest(sqlc.arg(digest_functions)::smallint[]) AS digest_function,
        unnest(sqlc.arg(hashes)::bytea[]) AS hash,
        unnest(sqlc.arg(size_bytes)::bigint[]) AS size_bytes
) AS input
JOIN file_paths fp
    ON fp.bep_instance_name_id = sqlc.arg(bep_instance_name_id)::bigint
    AND fp.path = input.file_path
JOIN digests d
    ON d.rev2_instance_name = input.rev2_instance_name
    AND d.digest_function = input.digest_function
    AND d.hash = input.hash
    AND d.size_bytes = input.size_bytes
JOIN files f
    ON f.file_path_id = fp.id
    AND f.digest_id = d.id;
