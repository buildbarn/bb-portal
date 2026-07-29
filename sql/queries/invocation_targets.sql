-- name: CreateInvocationTargetsBulk :exec
INSERT INTO invocation_targets (
    bazel_invocation_invocation_targets,
    target_invocation_targets,
    invocation_target_configuration,
    success,
    tags,
    failure_message,
    abort_reason
)
SELECT
    sqlc.arg(bazel_invocation_id),
    input.target_id,
    cfg.id,
    input.success,
    NULLIF(input.tags, '')::jsonb,
    NULLIF(input.failure_message, ''),
    input.abort_reason
FROM (
    SELECT
        unnest(sqlc.arg(target_ids)::bigint[]) AS target_id,
        unnest(sqlc.arg(configuration_ids)::text[]) AS configuration_external_id,
        unnest(sqlc.arg(successes)::boolean[]) AS success,
        unnest(sqlc.arg(tags_list)::text[]) AS tags,
        unnest(sqlc.arg(failure_messages)::text[]) AS failure_message,
        unnest(sqlc.arg(abort_reasons)::text[]) AS abort_reason
) AS input
JOIN configurations cfg
  ON cfg.bazel_invocation_id = sqlc.arg(bazel_invocation_id)
  AND cfg.configuration_id = input.configuration_external_id;
