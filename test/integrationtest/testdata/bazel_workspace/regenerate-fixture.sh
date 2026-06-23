#!/usr/bin/env bash
# Build the fixture workspace inside a Linux container and capture its BEP
# output as the integration test fixture. Run this when you change the
# workspace sources; the resulting fixture is checked in.
set -euo pipefail

WORKSPACE_DIR="$(cd "$(dirname "$0")" && pwd)"
FIXTURE_PATH="$(cd "$WORKSPACE_DIR/.." && pwd)/bepfiles/artifacts_endtoend.bep.ndjson"
IMAGE="gcr.io/bazel-public/bazel:latest"
TAR_FILE="${TMPDIR:-/tmp}/artifacts-fixture-workspace.tar"

cd "$WORKSPACE_DIR"
tar c --exclude='./bazel-*' --exclude='./regenerate-fixture.sh' . > "$TAR_FILE"

# Build the workspace, write the BEP to stdout (file descriptor 3 inside the
# container, captured to a host file via the trick at the bottom).
MSYS_NO_PATHCONV=1 docker run --rm -i \
  --entrypoint /usr/bin/bash \
  "$IMAGE" -c '
    set -e
    mkdir -p /home/ubuntu/work
    cd /home/ubuntu/work
    tar x
    bazel build //:hello \
      --build_event_json_file=/tmp/bep.ndjson \
      --noshow_progress 1>&2
    cat /tmp/bep.ndjson
  ' < "$TAR_FILE" > "$FIXTURE_PATH"

echo "Wrote $(wc -l < "$FIXTURE_PATH") events to $FIXTURE_PATH" >&2
