SHELL = bash -o pipefail

.PHONY: default
default: download lint test

.PHONY: download
download:
	go mod download

.PHONY: lint
lint:
	bazel mod tidy
	bazel run //:gazelle
	bazel run @com_github_bazelbuild_buildtools//:buildifier
	bazel run @cc_mvdan_gofumpt//:gofumpt -- -w -extra $(CURDIR)
	bazel run @org_golang_x_lint//golint -- -set_exit_status $(CURDIR)/...
	bazel test //...

.PHONY: lint-fix
lint-fix:
	golangci-lint --timeout 10m run --fix ./...

.PHONY: test
test:
	go test ./...

.PHONY: generate-bazel
generate-bazel:
	go generate ./third_party/bazel/...
