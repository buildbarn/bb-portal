## Purpose

There is no official generated Go code for the Build Event Stream protobuf definition. This Bazel project uses the
original protobuf definitions to generate Go code.

## Upgrading

From time to time, the dependency on Bazel should be upgraded, in order to keep up with the latest changes in the Build
Event Protocol (BEP). This can be done by editing the version of `com_github_bazelbuild_bazel` referenced at the bottom
of the `WORKSPACE` file.

NOTE: Dependencies of the Build Event Stream protobuf definition are managed manually. If there are some changes with
dependencies, the setup needs to be updated (e.g. a new import in the Build Event Stream protobuf definition is added).
