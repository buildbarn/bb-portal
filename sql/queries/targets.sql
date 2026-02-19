-- name: FindTargets :many
SELECT t.id, t.label, t.aspect, t.target_kind
FROM (
    SELECT 
        sqlc.arg(instance_name_id)::bigint AS instance_name_id,
        unnest(sqlc.arg(labels)::text[]) AS label,
        unnest(sqlc.arg(aspects)::text[]) AS aspect,
        unnest(sqlc.arg(target_kinds)::text[]) AS target_kind
) AS input
JOIN targets t ON 
    t.instance_name_targets = input.instance_name_id AND
    t.label = input.label AND
    t.aspect = input.aspect AND
    t.target_kind = input.target_kind;

-- name: CreateTargets :many
INSERT INTO targets (instance_name_targets, label, aspect, target_kind)
SELECT 
    instance_name_id, label, aspect, target_kind
FROM (
    SELECT 
        sqlc.arg(instance_name_id)::bigint AS instance_name_id,
        unnest(sqlc.arg(labels)::text[]) AS label,
        unnest(sqlc.arg(aspects)::text[]) AS aspect,
        unnest(sqlc.arg(target_kinds)::text[]) AS target_kind
) AS input
-- ORDER BY here is enforcing an insertion order for two reasons:
--   1. Prevent concurrent requests from deadlocking each other in case
-- they have overlapping insertions of targets in different orders.
--   2. Explicitly collate to normalize sort order between database
-- instantiations, otherwise golden file generation may have a different
-- order than what's used during the test.
ORDER BY
    label COLLATE "C",
    aspect COLLATE "C",
    target_kind COLLATE "C"
RETURNING id, label, aspect, target_kind;

-- name: DeleteUnusedTargetsFromPages :execrows
DELETE FROM targets
WHERE
    ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
    AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
    AND (
        NOT EXISTS (
            SELECT 1 FROM "invocation_targets" 
            WHERE "target_invocation_targets" = "targets"."id"
        )
    )
    AND (
        NOT EXISTS (
            SELECT 1 FROM "target_kind_mappings" 
            WHERE "target_id" = "targets"."id"
        )
    );
