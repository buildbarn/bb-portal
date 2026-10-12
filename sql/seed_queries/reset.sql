-- name: SeedTruncateEverythingCascade :exec
-- Deletes all seeded data and everything referencing it. This is an extremely
-- destructive operation and is only safe to run when bb_portal is stopped.
TRUNCATE TABLE instance_names, authenticated_users RESTART IDENTITY CASCADE;
