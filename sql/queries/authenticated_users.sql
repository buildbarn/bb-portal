-- name: CreateAuthenticatedUser :one
--
-- A function for creating authenticated
-- users. It returns the ID and a bool
-- `created`: `true` if the user was
-- created, `false` if the user already
-- existed. The query uses the system column
-- `xmax`, see docs for more information.
WITH new_row AS (
    INSERT INTO authenticated_users (
        user_uuid,
        external_id,
        display_name,
        user_info
    )
    VALUES (sqlc.arg(user_uuid), sqlc.arg(external_id), sqlc.arg(display_name), sqlc.arg(user_info))
    ON CONFLICT (external_id) DO UPDATE SET
        display_name = COALESCE(sqlc.arg(display_name), authenticated_users.display_name),
        user_info = COALESCE(sqlc.arg(user_info), authenticated_users.user_info)
    RETURNING id, (xmax = 0) AS created
)
SELECT id, created FROM new_row;
