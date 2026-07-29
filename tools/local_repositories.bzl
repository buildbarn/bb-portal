def _starlarkified_local_repository_impl(repository_ctx):
    relative_path = repository_ctx.attr.path
    workspace_root = repository_ctx.path(Label("//:MODULE.bazel")).dirname
    absolute_path = workspace_root
    for segment in relative_path.split("/"):
        absolute_path = absolute_path.get_child(segment)
    repository_ctx.symlink(absolute_path, ".")

starlarkified_local_repository = repository_rule(
    implementation = _starlarkified_local_repository_impl,
    attrs = {
        "path": attr.string(mandatory = True),
    },
)

_symlink = tag_class(
    attrs = {
        "name": attr.string(mandatory = True),
        "path": attr.string(mandatory = True),
    },
)

def _local_repositories_impl(ctx):
    for mod in ctx.modules:
        for symlink in mod.tags.symlink:
            starlarkified_local_repository(name = symlink.name, path = symlink.path)

local_repositories = module_extension(
    implementation = _local_repositories_impl,
    tag_classes = {
        "symlink": _symlink,
    },
)
