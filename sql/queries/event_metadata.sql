-- name: RecordEventMetadata :execresult
INSERT INTO event_metadata (
    sequence_number,
    event_received_at,
    event_hash,
    bazel_invocation_id
)
SELECT 
    sqlc.arg(sequence_number),
    sqlc.arg(event_received_at),
    sqlc.arg(event_hash),
    b.id
FROM bazel_invocations AS b
WHERE b.invocation_id = sqlc.arg(invocation_id)
  AND b.bep_completed = FALSE;
