-- name: GetOrCreateEventMetadata :one
WITH new_row AS (
    INSERT INTO event_metadata (
        bazel_invocation_id,
        handled,
        event_received_at,
        version
    ) VALUES (
        sqlc.arg(bazel_invocation_id), '\x', NOW(), 0
    )
    ON CONFLICT (bazel_invocation_id) DO NOTHING
    RETURNING id, handled, event_received_at, version
)
SELECT * FROM new_row
UNION ALL
SELECT id, handled, event_received_at, version
FROM event_metadata
WHERE bazel_invocation_id = sqlc.arg(bazel_invocation_id);

-- name: UpdateEventMetadata :one
UPDATE event_metadata
SET 
    handled = sqlc.arg(handled),
    event_received_at = sqlc.arg(event_received_at),
    version = version + 1
WHERE 
    id = sqlc.arg(id)
    AND version = sqlc.arg(version)
RETURNING version;
