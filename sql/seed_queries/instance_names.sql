-- name: SeedCreateInstanceNames :execrows
INSERT INTO instance_names (
    id,
    name
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(names)::text[]);

-- name: SeedCreateAuthenticatedUsers :execrows
INSERT INTO authenticated_users (
    id,
    user_uuid,
    external_id,
    display_name
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(user_uuids)::uuid[]),
    unnest(sqlc.arg(external_ids)::text[]),
    NULLIF(unnest(sqlc.arg(display_names)::text[]), '');

-- name: SeedAdvanceInstanceNamesSequence :exec
SELECT setval(
    pg_get_serial_sequence('instance_names', 'id'),
    COALESCE((SELECT MAX(id) FROM instance_names), 1)
) AS does_not_matter;

-- name: SeedAdvanceAuthenticatedUsersSequence :exec
SELECT setval(
    pg_get_serial_sequence('authenticated_users', 'id'),
    COALESCE((SELECT MAX(id) FROM authenticated_users), 1)
) AS does_not_matter;
