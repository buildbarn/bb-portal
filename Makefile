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
	cd frontend && npx graphql-codegen --config ./src/graphql/codegen.ts
	git restore pkg/proto/configuration/bb_portal/BUILD.bazel

.PHONY: lint-fix
lint-fix:
	golangci-lint --timeout 10m run --fix ./...

.PHONY: npxgen
npxgen:
	cd frontend && npx graphql-codegen --config ./src/graphql/codegen.ts

.PHONY: test
test:
	bazel test //...

.PHONY: update-schema
update-schema:
	go generate ./...
	bazel run //:gazelle
	git restore pkg/proto/configuration/bb_portal/BUILD.bazel

.PHONY: update-tests
update-tests:
	go test ./pkg/processing/ -snapshot-for-api-tests
	go test ./internal/graphql/ -update-golden
	go test ./pkg/summary/ -update-golden

.PHONY: generate-bazel
generate-bazel:
	go generate ./third_party/bazel/...
