-- name: CreateTestTargetsBulk :exec
INSERT INTO test_targets (target_id)
SELECT unnest(@target_ids::bigint[]) as target_id
ORDER BY target_id
ON CONFLICT (target_id) DO NOTHING;

-- name: DeleteOrphanedTestTargetsFromPages :execrows
DELETE FROM test_targets
WHERE
    ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
    AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
    AND NOT EXISTS (
        SELECT 1 
        FROM invocation_targets it
        JOIN test_summaries ts ON ts.invocation_target_test_summary = it.id
        WHERE it.target_invocation_targets = test_targets.target_id
    );
