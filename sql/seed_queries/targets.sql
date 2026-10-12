-- name: SeedCreateTargets :execrows
INSERT INTO targets (
    id,
    label,
    aspect,
    target_kind,
    instance_name_targets
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(labels)::text[]),
    unnest(sqlc.arg(aspects)::text[]),
    unnest(sqlc.arg(target_kinds)::text[]),
    unnest(sqlc.arg(instance_name_ids)::bigint[]);

-- name: SeedCreateTestTargets :execrows
INSERT INTO test_targets (
    id,
    target_id
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(target_ids)::bigint[]);

-- name: SeedAdvanceTargetsSequence :exec
SELECT setval(
    pg_get_serial_sequence('targets', 'id'),
    COALESCE((SELECT MAX(id) FROM targets), 1)
) AS does_not_matter;

-- name: SeedAdvanceTestTargetsSequence :exec
SELECT setval(
    pg_get_serial_sequence('test_targets', 'id'),
    COALESCE((SELECT MAX(id) FROM test_targets), 1)
) AS does_not_matter;
