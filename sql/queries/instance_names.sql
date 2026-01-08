-- name: CreateInstanceName :one
--
-- An idempotent function for creating an instance name. If the instance
-- name already exists, it will return the existing id.
WITH new_row AS (
    INSERT INTO instance_names (name) VALUES ($1)
    ON CONFLICT (name) DO NOTHING
    RETURNING id
)
SELECT id FROM new_row
UNION ALL
SELECT id FROM instance_names WHERE name = $1
LIMIT 1;
