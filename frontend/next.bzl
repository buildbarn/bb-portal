# Copyright 2024 The Bazel examples and tutorials Authors & Contributors.
# All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Wrapper macro to make the next.js rules more ergonomic."""

load("@aspect_rules_js//js:defs.bzl", "js_run_binary", "js_run_devserver")
load("@npm//frontend:next/package_json.bzl", next_bin = "bin")

def next(
        name,
        srcs,
        data,
        next_js_binary,
        next_bin,
        next_build_out = ".next",
        **kwargs):
    """Generates Next.js targets build, dev & start targets.

    `{name}`       - a js_run_binary build target that runs `next build`
    `{name}_dev`   - a js_run_devserver binary target that runs `next dev`
    `{name}_start` - a js_run_devserver binary target that runs `next start`

    Use this macro in the BUILD file at the root of a next app where the `next.config.js` file is
    located.

    For example, a target such as

    ```
    next(
        name = "next",
        srcs = [
            "//next.js/pages",
            "//next.js/public",
            "//next.js/styles",
        ],
        data = [
            "//next.js:node_modules/next",
            "//next.js:node_modules/react-dom",
            "//next.js:node_modules/react",
            "//next.js:node_modules/typescript",
            "next.config.js",
            "package.json",
        ],
        next_bin = "../../node_modules/.bin/next",
        next_js_binary = "//:next_js_binary",
    )
    ```

    in an `next.js/BUILD.bazel` file where the root `BUILD.bazel` file has
    next linked to `node_modules` and the `next_js_binary`:

    ```
    load("@npm//:defs.bzl", "npm_link_all_packages")
    load("@npm//:next/package_json.bzl", next_bin = "bin")

    npm_link_all_packages(name = "node_modules")

    next_bin.next_binary(
        name = "next_js_binary",
        visibility = ["//visibility:public"],
    )
    ```

    will create the targets:

    ```
    //next.js:next
    //next.js:next_dev
    //next.js:next_start
    ```

    To build the above next app, equivalent to running
    `next build` outside Bazel, run,

    ```
    bazel build //next.js:next
    ```

    To run the development server in watch mode with
    [ibazel](https://github.com/bazelbuild/bazel-watcher), equivalent to running
    `next dev` outside Bazel, run

    ```
    ibazel run //next.js:next_dev
    ```

    To run the production server in watch mode with
    [ibazel](https://github.com/bazelbuild/bazel-watcher), equivalent to running
    `next start` outside Bazel,

    ```
    ibazel run //next.js:next_start
    ```

    Args:
        name: The name of the build target.

        srcs: Source files to include in build & dev targets.
            Typically these are source files or transpiled source files in Next.js source folders
            such as `pages`, `public` & `styles`.

        data: Data files to include in all targets.
            These are typically npm packages required for the build & configuration files such as
            package.json and next.config.js.

        next_js_binary: The next js_binary. Used for the `build `target.

            Typically this is a js_binary target created using `bin` loaded from the `package_json.bzl`
            file of the npm package.

            See main docstring above for example usage.

        next_bin: The next bin command. Used for the `dev` and `start` targets.

            Typically the path to the next entry point from the current package. For example `./node_modules/.bin/next`,
            if next is linked to the current package, or a relative path such as `../node_modules/.bin/next`, if next is
            linked in the parent package.

            See main docstring above for example usage.

        next_build_out: The `next build` output directory. Defaults to `.next` which is the Next.js default output directory for the `build` command.

        **kwargs: Other attributes passed to all targets such as `tags`.
    """

    tags = kwargs.pop("tags", [])

    # `next build` creates an optimized bundle of the application
    # https://nextjs.org/docs/api-reference/cli#build
    js_run_binary(
        name = name,
        tool = next_js_binary,
        args = ["build", "--experimental-build-mode=compile"],
        srcs = srcs + data,
        mnemonic = "NextBuild",
        out_dirs = [next_build_out],
        chdir = native.package_name(),
        tags = tags,
        **kwargs
    )

    # `next dev` runs the application in development mode
    # https://nextjs.org/docs/api-reference/cli#development
    js_run_devserver(
        name = "{}_dev".format(name),
        command = next_bin,
        args = ["dev"],
        data = srcs + data,
        chdir = native.package_name(),
        tags = tags,
        **kwargs
    )

    # `next start` runs the application in production mode
    # https://nextjs.org/docs/api-reference/cli#production
    js_run_devserver(
        name = "{}_start".format(name),
        command = next_bin,
        args = ["start"],
        data = data + [name],
        chdir = native.package_name(),
        tags = tags,
        **kwargs
    )
