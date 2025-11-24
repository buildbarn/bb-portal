-- name: LockBazelInvocationCompletion :one
SELECT id, bep_completed
FROM bazel_invocations
WHERE id = sqlc.arg(id)
FOR SHARE;
