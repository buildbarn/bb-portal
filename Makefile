SHELL = bash -o pipefail

.PHONY: default
default: download lint test

.PHONY: download
download:
	go mod download

.PHONY: lint
lint:
	golangci-lint --timeout 10m run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint --timeout 10m run --fix ./...

.PHONY: test
test:
	go test ./...

.PHONY: generate-bazel
generate-bazel:
	go generate ./third_party/bazel/...
