load("@rules_oci//oci:defs.bzl", "oci_image", "oci_image_index", "oci_push")

def container_push(name, image, repository, component):
    oci_push(
        name = name,
        image = image,
        repository = repository + component,
        remote_tags = "@com_github_buildbarn_bb_storage//tools:stamped_tags",
        target_compatible_with = select({
            Label("@platforms//os:windows"): [Label("@platforms//:incompatible")],
            "//conditions:default": [],
        }),
    )
