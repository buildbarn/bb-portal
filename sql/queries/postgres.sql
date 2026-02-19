-- name: SelectPages :one
--
-- Returns the number of physical blocks of a table
SELECT 
    relpages
FROM pg_class 
WHERE relname = sqlc.arg(table_name);
