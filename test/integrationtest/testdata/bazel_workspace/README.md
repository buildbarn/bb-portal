# Bazel workspace for the artifact-listing integration test fixture

This is a minimal Bazel workspace whose purpose is to produce a deterministic
BEP file (with `NamedSetOfFiles` and `TargetCompleted` events) for the
`TestArtifactsEndToEnd` integration test.

The two source files (`greeting.cc` + `hello.cc`) build into two targets
(`:greeting` and `:hello`), each producing one default-output-group file.
The test fixture lives at `../bepfiles/artifacts_endtoend.bep.ndjson` and is
captured from a Linux container build of this workspace.

To regenerate the fixture:

```sh
test/integrationtest/testdata/bazel_workspace/regenerate-fixture.sh
```

The script runs Bazel inside `gcr.io/bazel-public/bazel:latest` so the
captured BEP is identical regardless of host OS. After regenerating, update
the `TestArtifactsEndToEnd` golden files:

```sh
bazel run //test/integrationtest:integrationtest_test -- --update-golden
```
