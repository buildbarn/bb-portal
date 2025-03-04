#!/bin/bash
set -eu

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT_DIR="$(realpath "$SCRIPT_DIR/../..")"
FRONTEND_DIR="${ROOT_DIR}/frontend"
PROTO_DIR="${ROOT_DIR}/frontend/src/proto/"
TS_PROTO_OPT="env=browser,outputServices=nice-grpc,outputServices=generic-definitions,useExactTypes=false,forceLong=string"

rm -rf ${FRONTEND_DIR}/src/lib/grpc-client/*

generate_grpc_client() {
  local proto_file=$1
  ${FRONTEND_DIR}/node_modules/.bin/grpc_tools_node_protoc \
    --plugin=protoc-gen-ts_proto="${FRONTEND_DIR}/node_modules/.bin/protoc-gen-ts_proto" \
    --ts_proto_out=${FRONTEND_DIR}/src/lib/grpc-client/ \
    --ts_proto_opt=${TS_PROTO_OPT} \
    --proto_path=${PROTO_DIR}/ \
    ${proto_file}
}

generate_grpc_client "${PROTO_DIR}/buildbarn/auth/auth.proto"
generate_grpc_client "${PROTO_DIR}/buildbarn/buildqueuestate/buildqueuestate.proto"
generate_grpc_client "${PROTO_DIR}/buildbarn/cas/cas.proto"
generate_grpc_client "${PROTO_DIR}/buildbarn/resourceusage/resourceusage.proto"
generate_grpc_client "${PROTO_DIR}/google/bytestream/bytestream.proto"
generate_grpc_client "${PROTO_DIR}/buildbarn/iscc/iscc.proto"
generate_grpc_client "${PROTO_DIR}/buildbarn/fsac/fsac.proto"
generate_grpc_client "${PROTO_DIR}/buildbarn/query/query.proto"
