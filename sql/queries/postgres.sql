-- name: SelectPages :one
--
-- Returns the number of physical blocks of a table
SELECT 
    relpages
FROM pg_class 
WHERE relname = sqlc.arg(table_name);

-- name: SelectForeignKeysWithoutIndexes :many
SELECT 
    c.conrelid::regclass::text AS table_name,
    c.conname AS foreign_key_name
FROM pg_constraint c
JOIN pg_namespace n ON n.oid = c.connamespace
WHERE c.contype = 'f'
  AND n.nspname = 'public'
  AND NOT EXISTS (
      SELECT 1
      FROM pg_index i
      WHERE i.indrelid = c.conrelid
        AND (string_to_array(i.indkey::text, ' '))[1] = (c.conkey)[1]::text
  )
ORDER BY table_name, foreign_key_name;
