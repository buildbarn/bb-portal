#!/bin/bash
set -euo pipefail

MYDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOTDIR="$(realpath "$MYDIR/..")"
BAZELDIR="${ROOTDIR}/_generate-via-bazel"
REPO="github.com/buildbarn/bb-portal"
REPO_ROOT="$(realpath "$ROOTDIR/../..")"
GEN_PATH="third_party/bazel/gen"
GEN_SOURCE="${BAZELDIR}/bazel-bin/gopath/src/${REPO}/${GEN_PATH}"
GEN_TARGET="${REPO_ROOT}/${GEN_PATH}"

create_gopath_via_bazel() {
	(
		cd "${BAZELDIR}"
		bazel build //:gopath
	)
}

extract_generated_code() {
	local from=$1
	local to=$2
	mkdir -p "$to"
	cp "${from}"/* "${to}/"
	chmod 0644 "${to}"/*
}

create_gopath_via_bazel
extract_generated_code "${GEN_SOURCE}/bes" "${GEN_TARGET}/bes"
extract_generated_code "${GEN_SOURCE}/bescore" "${GEN_TARGET}/bescore"
