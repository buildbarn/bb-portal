#!/usr/bin/env bash

set -eux

MYDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
cd "$MYDIR"

bazel run @@//protobuf:protobuf
bazel test //...