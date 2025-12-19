#!/bin/bash

set -eEx

if [ -n "$BUILD_WORKSPACE_DIRECTORY" ]; then
  cd "$BUILD_WORKSPACE_DIRECTORY"
fi

go generate
bazel run //:gazelle
bazel run //cmd/bb_export_schema > sql/migrations/schema.sql
bazel run //tools:sqlc -- generate
