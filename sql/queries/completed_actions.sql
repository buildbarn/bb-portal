-- name: DeleteOldCompletedActions :execrows
DELETE FROM completed_actions
WHERE ctid IN(
    SELECT ctid
    FROM completed_actions
    WHERE
      ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
      AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
      AND created_timestamp < sqlc.arg(cutoff_time)::timestamptz
    FOR UPDATE SKIP LOCKED
    LIMIT sqlc.arg(batch_limit)::bigint
);
