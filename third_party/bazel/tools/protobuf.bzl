def _get_proto_src_impl(ctx):
    proto_inputs = ctx.attr.proto[ProtoInfo].direct_sources
    srcs = [
        src
        for src in proto_inputs
        if src.path.endswith(ctx.attr.filename_suffix)
    ]
    if len(srcs) != 1:
        fail("Could not find a single source %s in %s", ctx.attr.filename_suffix, proto_inputs)
    return DefaultInfo(
        files = depset(srcs),
    )

get_proto_src = rule(
    implementation = _get_proto_src_impl,
    attrs = {
        "proto": attr.label(
            mandatory = True,
            providers = [ProtoInfo],
        ),
        "filename_suffix": attr.string(
            doc = "Suffix of the file to extract.",
            mandatory = True,
        ),
    },
)
