-- name: InsertMissingFilePaths :exec
INSERT INTO file_paths (bep_instance_name_id, path)
SELECT
    sqlc.arg(bep_instance_name_id)::bigint,
    input.path
FROM (
    SELECT
        unnest(sqlc.arg(paths)::text[]) AS path
) AS input
ORDER BY input.path
ON CONFLICT (bep_instance_name_id, path) DO NOTHING;

-- name: InsertMissingDigests :exec
INSERT INTO digests (rev2_instance_name, digest_function, hash, size_bytes)
SELECT
    input.instance_name,
    input.digest_function,
    input.hash,
    input.size_bytes
FROM (
    SELECT
        unnest(sqlc.arg(rev2_instance_names)::text[]) AS instance_name,
        unnest(sqlc.arg(digest_functions)::smallint[]) AS digest_function,
        unnest(sqlc.arg(hashes)::bytea[]) AS hash,
        unnest(sqlc.arg(size_bytes)::bigint[]) AS size_bytes
) AS input
ORDER BY input.hash
ON CONFLICT (rev2_instance_name, digest_function, hash, size_bytes) DO NOTHING;

-- name: InsertMissingFiles :exec
INSERT INTO files (
    digest_id,
    file_path_id
)
SELECT
    d.id,
    fp.id
FROM (
    SELECT
        unnest(sqlc.arg(file_paths)::text[]) AS file_path,
        unnest(sqlc.arg(rev2_instance_names)::text[]) AS rev2_instance_name,
        unnest(sqlc.arg(digest_functions)::smallint[]) digest_function,
        unnest(sqlc.arg(hashes)::bytea[]) AS hash,
        unnest(sqlc.arg(size_bytes)::bigint[]) AS size_bytes
) AS input
JOIN file_paths fp
    ON fp.bep_instance_name_id = sqlc.arg(bep_instance_name_id)::bigint
    AND fp.path = input.file_path
JOIN digests d
    ON d.rev2_instance_name = input.rev2_instance_name
    AND d.digest_function = input.digest_function
    AND d.hash = input.hash
    AND d.size_bytes = input.size_bytes
ORDER BY d.id, fp.id
ON CONFLICT (digest_id, file_path_id)
DO NOTHING;

-- name: DeleteUnusedFilePathsFromPages :execrows
DELETE FROM file_paths
WHERE ctid IN (
    SELECT ctid
    FROM file_paths
    WHERE
        ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
        AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
        AND (
            NOT EXISTS (
                SELECT 1 FROM "files"
                WHERE "file_path_id" = "file_paths"."id"
            )
        )
    FOR UPDATE SKIP LOCKED
    LIMIT sqlc.arg(batch_limit)::bigint
);

-- name: DeleteUnusedDigestFromPages :execrows
DELETE FROM digests
WHERE ctid IN (
    SELECT ctid
    FROM digests
    WHERE
        ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
        AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
        AND (
            NOT EXISTS (
                SELECT 1 FROM "files"
                WHERE "digest_id" = "digests"."id"
            )
        )
    FOR UPDATE SKIP LOCKED
    LIMIT sqlc.arg(batch_limit)::bigint
);

-- name: DeleteUnusedFilesFromPages :execrows
DELETE FROM files
WHERE ctid IN (
    SELECT ctid
    FROM files
    WHERE
        ctid >= format('(%s,0)', sqlc.arg(from_page)::bigint)::tid
        AND ctid < format('(%s,0)', sqlc.arg(from_page)::bigint + sqlc.arg(pages)::bigint)::tid
        AND (
            NOT EXISTS (
                SELECT 1 FROM "actions"
                WHERE "stdout_file_id" = "files"."id"
                  OR "stderr_file_id" = "files"."id"
            )
        )
        AND (
            NOT EXISTS (
                SELECT 1 FROM "build_tool_logs"
                WHERE "file_id" = "files"."id"
            )
        )
        AND (
            NOT EXISTS (
                SELECT 1 FROM "test_action_outputs"
                WHERE "file_id" = "files"."id"
            )
        )
    FOR UPDATE SKIP LOCKED
    LIMIT sqlc.arg(batch_limit)::bigint
);
