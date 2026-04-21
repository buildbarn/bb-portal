-- name: UpdateBuildTimestampFromInvocation :exec
UPDATE builds
SET timestamp = bi.started_at
FROM bazel_invocations bi
WHERE bi.id = sqlc.arg(invocation_id)
  AND builds.id = bi.build_invocations
  AND bi.started_at < builds.timestamp;
