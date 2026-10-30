-- name: SeedCreateBazelInvocations :execrows
INSERT INTO bazel_invocations (
    id,
    invocation_id,
    created_timestamp,
    started_at,
    ended_at,
    bep_completed,
    username,
    hostname,
    exit_code_name,
    exit_code_code,
    authenticated_user_bazel_invocations,
    instance_name_bazel_invocations
)
SELECT
    s.id,
    s.invocation_id,
    s.created_timestamp,
    s.started_at,
    NULLIF(s.ended_at, '0001-01-01 00:00:00+00'::timestamptz), -- zero time is a sentinel for NULL
    s.bep_completed,
    NULLIF(s.username, ''), -- empty string is a sentinel for NULL
    NULLIF(s.hostname, ''),
    NULLIF(s.exit_code_name, ''),
    CASE WHEN s.exit_code_code = -1 THEN NULL ELSE s.exit_code_code END, -- -1 is a sentinel for NULL
    CASE WHEN s.authenticated_user_id = 0 THEN NULL ELSE s.authenticated_user_id END, -- 0 is a sentinel for NULL
    s.instance_name_id
FROM (
    SELECT
        unnest(sqlc.arg(ids)::bigint[]) AS id,
        unnest(sqlc.arg(invocation_ids)::uuid[]) AS invocation_id,
        unnest(sqlc.arg(created_timestamps)::timestamptz[]) AS created_timestamp,
        unnest(sqlc.arg(started_ats)::timestamptz[]) AS started_at,
        unnest(sqlc.arg(ended_ats)::timestamptz[]) AS ended_at,
        unnest(sqlc.arg(bep_completeds)::boolean[]) AS bep_completed,
        unnest(sqlc.arg(usernames)::text[]) AS username,
        unnest(sqlc.arg(hostnames)::text[]) AS hostname,
        unnest(sqlc.arg(exit_code_names)::text[]) AS exit_code_name,
        unnest(sqlc.arg(exit_code_codes)::int[]) AS exit_code_code,
        unnest(sqlc.arg(authenticated_user_ids)::bigint[]) AS authenticated_user_id,
        unnest(sqlc.arg(instance_name_ids)::bigint[]) AS instance_name_id
) AS s;

-- name: SeedCreateSourceControls :execrows
INSERT INTO source_controls (
    id,
    repo,
    repo_url,
    ref,
    ref_url,
    commit,
    commit_url,
    bazel_invocation_source_control
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(repos)::text[]),
    unnest(sqlc.arg(repo_urls)::text[]),
    unnest(sqlc.arg(refs)::text[]),
    unnest(sqlc.arg(ref_urls)::text[]),
    unnest(sqlc.arg(commits)::text[]),
    unnest(sqlc.arg(commit_urls)::text[]),
    unnest(sqlc.arg(bazel_invocation_ids)::bigint[]);

-- name: SeedCreateInvocationTags :execrows
INSERT INTO invocation_tags (
    id,
    key,
    value,
    bazel_invocation_id
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(keys)::text[]),
    unnest(sqlc.arg(values)::text[]),
    unnest(sqlc.arg(bazel_invocation_ids)::bigint[]);

-- name: SeedCreateConnectionMetadata :execrows
INSERT INTO connection_metadata (
    id,
    connection_last_open_at,
    bazel_invocation_connection_metadata
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(connection_last_open_ats)::timestamptz[]),
    unnest(sqlc.arg(bazel_invocation_ids)::bigint[]);

-- name: SeedCreateEventMetadata :execrows
INSERT INTO event_metadata (
    id,
    handled,
    event_received_at,
    version,
    bazel_invocation_id
)
SELECT
    unnest(sqlc.arg(ids)::bigint[]),
    unnest(sqlc.arg(handleds)::bytea[]),
    unnest(sqlc.arg(event_received_ats)::timestamptz[]),
    unnest(sqlc.arg(versions)::bigint[]),
    unnest(sqlc.arg(bazel_invocation_ids)::bigint[]);

-- name: SeedAdvanceBazelInvocationsSequence :exec
SELECT setval(
    pg_get_serial_sequence('bazel_invocations', 'id'),
    COALESCE((SELECT MAX(id) FROM bazel_invocations), 1)
) AS does_not_matter;

-- name: SeedAdvanceSourceControlsSequence :exec
SELECT setval(
    pg_get_serial_sequence('source_controls', 'id'),
    COALESCE((SELECT MAX(id) FROM source_controls), 1)
) AS does_not_matter;

-- name: SeedAdvanceInvocationTagsSequence :exec
SELECT setval(
    pg_get_serial_sequence('invocation_tags', 'id'),
    COALESCE((SELECT MAX(id) FROM invocation_tags), 1)
) AS does_not_matter;

-- name: SeedAdvanceConnectionMetadataSequence :exec
SELECT setval(
    pg_get_serial_sequence('connection_metadata', 'id'),
    COALESCE((SELECT MAX(id) FROM connection_metadata), 1)
) AS does_not_matter;

-- name: SeedAdvanceEventMetadataSequence :exec
SELECT setval(
    pg_get_serial_sequence('event_metadata', 'id'),
    COALESCE((SELECT MAX(id) FROM event_metadata), 1)
) AS does_not_matter;
