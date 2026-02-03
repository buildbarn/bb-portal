"""Rule to unpack a .tar.xz and make it available as a directory."""

def _unpack_impl(ctx):
    input_tar = ctx.file.src
    output_dir = ctx.actions.declare_directory(ctx.attr.name + ".extracted")
    args = ctx.actions.args()
    args.add_all([
        ctx.file._xz.path,
        ctx.file._tar.path,
        input_tar.path,
        output_dir.path,
    ])
    ctx.actions.run_shell(
        inputs = [input_tar, ctx.file._xz, ctx.file._tar],
        outputs = [output_dir],
        arguments = [args],
        command = """
        xz=$1
        tar=$2
        input=$3
        output=$4
        "$tar" -I "$xz" -xf "$input" -C "$output"
        """
    )
    return [DefaultInfo(files = depset([output_dir]))]


unpack = rule(
    implementation = _unpack_impl,
    attrs = {
        "src": attr.label(allow_single_file = True, mandatory = True),
        "_tar": attr.label(default = "@ape//ape:tar", allow_single_file = True, executable = True, cfg = 'exec'),
        "_xz": attr.label(default = "@ape//ape:xz", allow_single_file = True, executable = True, cfg = 'exec'),
    },
)
