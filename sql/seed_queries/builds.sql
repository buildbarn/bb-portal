-- name: SeedCreateBuilds :execrows
INSERT INTO builds (
    id,
    build_uuid,
    timestamp,
    instance_name_builds
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(build_uuids)::uuid[]),
    unnest(sqlc.arg(timestamps)::timestamptz[]),
    unnest(sqlc.arg(instance_name_ids)::bigint[]);

-- name: SeedAssignInvocationsToBuild :execrows
-- Associates invocations with a build by setting their build foreign key.
UPDATE bazel_invocations
SET build_invocations = sqlc.arg(build_id)::bigint
WHERE id = ANY(sqlc.arg(invocation_ids)::bigint[]);

-- name: SeedCreateBuildTags :execrows
INSERT INTO build_tags (
    id,
    key,
    value,
    build_id
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(keys)::text[]),
    unnest(sqlc.arg(values)::text[]),
    unnest(sqlc.arg(build_ids)::bigint[]);

-- name: SeedAdvanceBuildsSequence :exec
SELECT setval(
    pg_get_serial_sequence('builds', 'id'),
    COALESCE((SELECT MAX(id) FROM builds), 1)
) AS does_not_matter;

-- name: SeedAdvanceBuildTagsSequence :exec
SELECT setval(
    pg_get_serial_sequence('build_tags', 'id'),
    COALESCE((SELECT MAX(id) FROM build_tags), 1)
) AS does_not_matter;
