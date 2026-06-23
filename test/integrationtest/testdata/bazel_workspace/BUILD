load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library")

cc_library(
    name = "greeting",
    srcs = ["greeting.cc"],
    hdrs = ["greeting.h"],
)

cc_binary(
    name = "hello",
    srcs = ["hello.cc"],
    deps = [":greeting"],
)
