"""Rule to unpack a .tar.xz and make it available as a directory."""

def _unpack_impl(ctx):
    input_tar = ctx.file.src
    output_dir = ctx.actions.declare_directory(ctx.attr.name + ".extracted")
    args = ctx.actions.args()
    args.add("-I", ctx.executable._xz)
    args.add("--touch")
    args.add("-xf", input_tar)
    args.add("-C", output_dir.path)
    ctx.actions.run(
        executable = ctx.executable._tar,
        inputs = [input_tar],
        tools = [ctx.executable._xz],
        outputs = [output_dir],
        arguments = [args],
    )
    return [DefaultInfo(files = depset([output_dir]))]

unpack = rule(
    implementation = _unpack_impl,
    attrs = {
        "src": attr.label(allow_single_file = True, mandatory = True),
        "_tar": attr.label(default = "@ape//ape:tar", executable = True, cfg = "exec"),
        "_xz": attr.label(default = "@ape//ape:xz", executable = True, cfg = "exec"),
    },
)
