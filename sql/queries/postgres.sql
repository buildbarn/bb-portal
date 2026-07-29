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

-- name: SelectRedundantIndexes :many
WITH index_info AS (
    SELECT
        c.relname AS table_name,
        i.indexrelid,
        ic.relname AS index_name,
        i.indrelid,
        i.indisunique,
        string_to_array(i.indkey::text, ' ') AS index_cols
    FROM pg_index i
    JOIN pg_class ic ON i.indexrelid = ic.oid
    JOIN pg_class c ON i.indrelid = c.oid
    JOIN pg_namespace n ON c.relnamespace = n.oid
    WHERE n.nspname = 'public'
)
SELECT
    i1.table_name,
    i1.index_name AS redundant_index,
    i2.index_name AS covering_index
FROM index_info i1
JOIN index_info i2 
    ON i1.indrelid = i2.indrelid 
    AND i1.indexrelid != i2.indexrelid
WHERE 
    i1.indisunique = false
    AND i2.index_cols[1:array_length(i1.index_cols, 1)] = i1.index_cols
ORDER BY 
    i1.table_name, 
    i1.index_name;
