-- name: CreateTargetKindMappingsBulk :exec
INSERT INTO target_kind_mappings (
    bazel_invocation_id,
    target_id
)
SELECT
    sqlc.arg(bazel_invocation_id),
    target_id
FROM (
    SELECT unnest(sqlc.arg(target_ids)::bigint[]) AS target_id
) AS input;

-- name: FindMappedTargets :many
SELECT
    t.label,
    t.aspect,
    m.target_id
FROM (
    SELECT
        unnest(sqlc.arg(labels)::text[]) AS label,
        unnest(sqlc.arg(aspects)::text[]) AS aspect
) AS input
JOIN targets t ON
    t.label = input.label AND
    t.aspect = input.aspect
JOIN target_kind_mappings m ON
    m.target_id = t.id AND
    m.bazel_invocation_id = sqlc.arg(bazel_invocation_id);

-- name: DeleteTargetKindMappingsFromPages :execrows
DELETE FROM target_kind_mappings
WHERE ctid IN (
    SELECT m.ctid FROM target_kind_mappings m
    JOIN bazel_invocations AS i ON i.id = m.bazel_invocation_id
    WHERE
        m.bazel_invocation_id = i.id
        AND m.ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
        AND m.ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
        AND i.bep_completed = true
    FOR UPDATE SKIP LOCKED
    LIMIT sqlc.arg(batch_limit)::bigint
);
