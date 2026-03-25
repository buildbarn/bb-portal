load("@bazel_lib//lib:copy_to_bin.bzl", _copy_to_bin = "copy_to_bin")

def _vite_build_rule_impl(ctx):
    out_dir = ctx.actions.declare_directory(ctx.label.name)

    args = ctx.actions.args()
    args.add("build")
    num_up = 3 + len(ctx.attr.chdir.split("/"))
    args.add_all("--outDir", [out_dir], format_each = ("../" * num_up) + "%s", expand_directories = False)
    ctx.actions.run(
        executable = ctx.executable.vite_binary,
        arguments = [args],
        inputs = ctx.files.srcs,
        outputs = [out_dir],
        env = {
            "BAZEL_BINDIR": "/".join(ctx.executable.vite_binary.path.split("/")[0:3]),
            "JS_BINARY__CHDIR": ctx.attr.chdir,
            "JS_BINARY__SILENT_ON_SUCCESS": "1",
        } | ctx.attr.env,
    )
    return DefaultInfo(files = depset([out_dir]))

vite_build_rule = rule(
    implementation = _vite_build_rule_impl,
    attrs = {
        "srcs": attr.label_list(mandatory = True, allow_files = True, cfg = "exec"),
        "vite_binary": attr.label(mandatory = True, executable = True, cfg = "exec"),
        "env": attr.string_dict(),
        "chdir": attr.string(default = ""),
    },
    toolchains = [
        # Even if not used, make sure the analysis selects an execution platform
        # that we are likely to be able to build _vite_binary for.
        "@rules_nodejs//nodejs:runtime_toolchain_type",
    ],
)

def _vite_build_macro_impl(name, srcs, tags, **kwargs):
    copy_to_bin_name = "{}_copy_srcs_to_bin".format(name)
    _copy_to_bin(
        name = copy_to_bin_name,
        srcs = srcs,
        # Always tag the target manual since we should only build it when the final target is built.
        tags = (tags or []) + ["manual"],
    )
    vite_build_rule(
        name = name,
        srcs = [":{}".format(copy_to_bin_name)],
        tags = tags,
        **kwargs,
    )

vite_build = macro(
    doc = "Builds a project with vite.",
    implementation = _vite_build_macro_impl,
    inherit_attrs = vite_build_rule,
)
