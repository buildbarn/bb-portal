-- name: SeedCreateInvocationTargets :execrows
INSERT INTO invocation_targets (
    id,
    success,
    abort_reason,
    bazel_invocation_invocation_targets,
    target_invocation_targets
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(successes)::boolean[]),
    unnest(sqlc.arg(abort_reasons)::text[]),
    unnest(sqlc.arg(bazel_invocation_ids)::bigint[]),
    unnest(sqlc.arg(target_ids)::bigint[]);

-- name: SeedCreateTestSummaries :execrows
INSERT INTO test_summaries (
    id,
    overall_status,
    total_run_count,
    run_count,
    attempt_count,
    shard_count,
    total_num_cached,
    total_run_duration_in_ms,
    invocation_target_test_summary
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(overall_statuses)::text[]),
    unnest(sqlc.arg(total_run_counts)::int[]),
    unnest(sqlc.arg(run_counts)::int[]),
    unnest(sqlc.arg(attempt_counts)::int[]),
    unnest(sqlc.arg(shard_counts)::int[]),
    unnest(sqlc.arg(total_num_cacheds)::int[]),
    unnest(sqlc.arg(total_run_duration_in_mss)::bigint[]),
    unnest(sqlc.arg(invocation_target_ids)::bigint[]);

-- name: SeedCreateTestResults :execrows
INSERT INTO test_results (
    id,
    run,
    shard,
    attempt,
    status,
    cached_locally,
    test_attempt_duration_in_ms,
    strategy,
    cached_remotely,
    test_summary_test_results
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(runs)::int[]),
    unnest(sqlc.arg(shards)::int[]),
    unnest(sqlc.arg(attempts)::int[]),
    unnest(sqlc.arg(statuses)::text[]),
    unnest(sqlc.arg(cached_locallys)::boolean[]),
    unnest(sqlc.arg(test_attempt_duration_in_mss)::bigint[]),
    unnest(sqlc.arg(strategies)::text[]),
    unnest(sqlc.arg(cached_remotelys)::boolean[]),
    unnest(sqlc.arg(test_summary_ids)::bigint[]);

-- name: SeedAdvanceInvocationTargetsSequence :exec
SELECT setval(
    pg_get_serial_sequence('invocation_targets', 'id'),
    COALESCE((SELECT MAX(id) FROM invocation_targets), 1)
) AS does_not_matter;

-- name: SeedAdvanceTestSummariesSequence :exec
SELECT setval(
    pg_get_serial_sequence('test_summaries', 'id'),
    COALESCE((SELECT MAX(id) FROM test_summaries), 1)
) AS does_not_matter;

-- name: SeedAdvanceTestResultsSequence :exec
SELECT setval(
    pg_get_serial_sequence('test_results', 'id'),
    COALESCE((SELECT MAX(id) FROM test_results), 1)
) AS does_not_matter;
