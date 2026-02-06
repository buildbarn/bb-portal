-- name: LockBazelInvocationCompletion :one
SELECT id, bep_completed
FROM bazel_invocations
WHERE id = sqlc.arg(id)
FOR SHARE;

-- name: CreateBazelInvocation :one
--
-- An idempotent function for creating bazel invocations. If the
-- invocation already exists, it will return the existing id.
WITH new_row AS (
    INSERT INTO bazel_invocations (
        invocation_id,
        instance_name_bazel_invocations,
        authenticated_user_bazel_invocations
    )
    VALUES (sqlc.arg(invocation_id), sqlc.arg(instance_name_id), sqlc.arg(authenticated_user_id))
    ON CONFLICT (invocation_id) DO NOTHING
    RETURNING id
)
SELECT id FROM new_row
UNION ALL
SELECT id FROM bazel_invocations
WHERE invocation_id = sqlc.arg(invocation_id)
  AND instance_name_bazel_invocations = sqlc.arg(instance_name_id)
  AND authenticated_user_bazel_invocations IS NOT DISTINCT FROM sqlc.arg(authenticated_user_id)
  AND bep_completed = false
LIMIT 1;

-- name: LockStaleInvocationsFromPages :execrows
UPDATE bazel_invocations
SET bep_completed = true
WHERE
    ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
    AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
    AND bep_completed = false
    AND EXISTS (
        SELECT 1 FROM event_metadata em
        WHERE em.bazel_invocation_id = bazel_invocations.id
        AND em.event_received_at <= sqlc.arg(cutoff_time)::timestamptz
    );

-- name: UpdateCompletedInvocationWithEndTimeFromEventMetadata :execrows
UPDATE bazel_invocations bi
SET
  ended_at = em.event_received_at
FROM event_metadata em
WHERE bi.id = em.bazel_invocation_id
  AND bi.bep_completed = true
  AND bi.ended_at IS NULL;

-- name: DeleteOldInvocationsFromPages :execrows
DELETE FROM bazel_invocations
WHERE
    ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
    AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
    AND ended_at < sqlc.arg(cutoff_time)::timestamptz
    AND bep_completed = true;
