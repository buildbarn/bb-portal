#!/bin/bash

# --- begin runfiles.bash initialization v3 ---
# Copy-pasted from the Bazel Bash runfiles library v3.
set -uo pipefail; set +e; f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" 2>/dev/null || \
  source "$0.runfiles/$f" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
  { echo>&2 "ERROR: cannot find $f"; exit 1; }; f=; set -e
# --- end runfiles.bash initialization v3 ---

# Order must match order in BUILD.bazel
go_bazel_path="$1"
gofumpt_bazel_path="$2"
clang_format_bazel_path="$3"
gazelle_bazel_path="$4"
sqlc_bazel_path="$5"
bb_export_schema_bazel_path="$6"

# Resolve them to absolute paths
go="$(rlocation "$go_bazel_path")"
gofumpt="$(rlocation "$gofumpt_bazel_path")"
clang_format="$(rlocation "$clang_format_bazel_path")"
gazelle="$(rlocation "$gazelle_bazel_path")"
sqlc="$(rlocation "$sqlc_bazel_path")"
bb_export_schema="$(rlocation "$bb_export_schema_bazel_path")"

# List of variable names to validate
cmds=(go gofumpt clang_format gazelle sqlc bb_export_schema)
for cmd in "${cmds[@]}"; do
  cmd_path="${!cmd}"
  
  # Check if the path is empty OR the file is not executable
  if [[ -z "$cmd_path" ]] || [[ ! -x "$cmd_path" ]]; then
    echo "Error: Missing or non-executable dependency binary for '$cmd'." >&2
    echo "Path resolved to: ${cmd_path:-[EMPTY]}" >&2
    exit 1
  fi
done


# Start in the root directory
cd "$BUILD_WORKSPACE_DIRECTORY"

# Get the go module name
go_module_name=$($go list -m)

# Go dependencies
find bazel-bin/ -path "*${go_module_name}*" -name '*.pb.go' -delete || true
bazel build $(bazel query --output=label 'kind("go_proto_library", //...)')
find bazel-bin/ -path "*${go_module_name}*" -name '*.pb.go' | while read f; do
  cat "$f" > $(echo "$f" | sed -e "s|.*/${go_module_name}/||")
done

#$go get -d -u ./... || true
$go mod tidy || true

# Generate database files
$go generate
if ! $bb_export_schema > sql/migrations/schema.sql; then  
  echo "Schema export failed, this may be due BUILD.bazel files being outdated."
  echo "Please run 'bazel run //:gazelle' and try again."
  exit 1
fi
$sqlc generate

# Gazelle
find . -name '*.pb.go' -not -path './pkg/proto/bazelbuild/*' -delete
rm -f $(find . -name '*.proto' -not -path './pkg/proto/bazelbuild/*' -not -path './third_party/*' | sed -e 's/[^/]*$/BUILD.bazel/')
$gazelle

# bzlmod
bazel mod tidy

# Go
$gofumpt -w -extra "$(pwd)"

# Protobuf
find . -name '*.proto' -not -path './frontend/*' -exec "$clang_format" -i {} +

# Generated .pb.go files
find bazel-bin/ -path "*${go_module_name}*" -name '*.pb.go' -delete || true
bazel build --output_groups=go_generated_srcs $(bazel query --output=label 'kind("go_proto_library", //...)')
third_party/bazel/download_protofiles.sh
find bazel-bin/ -path "*${go_module_name}*" -name '*.pb.go' | while read f; do
  cat $f > $(echo $f | sed -e 's|^bazel-bin/||' -e 's|/[^/]*_go_proto_/.*/|/|')
done

# Files embedded into Go binaries
if git grep -q '^[[:space:]]*//go:embed '; then
  bazel build $(git grep '^[[:space:]]*//go:embed ' | sed -e 's|\(.*\)/.*//go:embed |//\1:|; s|"||g; s| .*||' | sort -u)
  git grep '^[[:space:]]*//go:embed ' | sed -e 's|\(.*\)/.*//go:embed |\1/|' | while read o; do
    if [ -e "bazel-bin/$o" ]; then
      rm -rf "$o"
      cp -r "bazel-bin/$o" "$o"
      find "$o" -type f -exec chmod -x {} +
    fi
  done
fi
