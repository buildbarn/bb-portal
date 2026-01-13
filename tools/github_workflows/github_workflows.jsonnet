// Encapsulates the complex multi-line installation scripts
local install_bazel(os_param, arch_param="x86_64") =
  local is_matrix = (os_param == "matrix");
  local bazel_os = if is_matrix then "${{matrix.host.bazel_os}}" else os_param;
  local bazel_arch = if is_matrix then "${{matrix.host.bazel_arch}}" else arch_param;

  {
    name: "Installing Bazel",
    shell: "bash",
    // We assume the host is always *nix, so we output to standard '~/bazel'
    run: |||
      v=$(cat .bazelversion)
      mkdir -p ~/bin
      curl -L https://github.com/bazelbuild/bazel/releases/download/${v}/bazel-${v}-%s-%s > ~/bin/bazel
      chmod +x ~/bin/bazel
      echo ~/bin >> ${GITHUB_PATH}
    ||| % [bazel_os, bazel_arch],
  };

// Build and test for a specific platform
local bazel_step(platform, if_cond) = {
  name: platform + ": build and test",
  'if': if_cond,
  run: "bazel test --test_output=errors --platforms=@rules_go//go/toolchain:%s_cgo //..." % platform,
};

// Define the list of platforms that need to be cross-compiled
local cross_targets = [
  "darwin_amd64",
  "darwin_arm64",
  "freebsd_amd64",
  "windows_amd64"
];

// A single step that builds for all target cross platforms
local cross_build_step = {
  name: "Cross-platform builds",
  'if': "matrix.host.cross_compile",
  local platforms = std.join(",", [
    "@rules_go//go/toolchain:%s_cgo" % p for p in cross_targets
  ]),
  run: "bazel build --platforms=%s //..." % platforms,
};

// Shared matrix definition
local build_matrix = {
  host: [
    {
      runs_on: "ubuntu-24.04",
      bazel_os: "linux",
      bazel_arch: "x86_64",
      platform_name: "linux_amd64",
      cross_compile: false,
    },
    {
      runs_on: "ubuntu-24.04-arm",
      bazel_os: "linux",
      bazel_arch: "arm64",
      platform_name: "linux_arm64",
      cross_compile: false,
    },
    {
      runs_on: "ubuntu-24.04-arm",
      bazel_os: "linux",
      bazel_arch: "arm64",
      platform_name: "cross_compiler",
      cross_compile: true,
    },
  ],
};

local checkout_step = { name: "Check out source code", uses: "actions/checkout@v4" };
local docker_credentials_step = {
  name: "Install Docker credentials",
  run: "echo \"${GITHUB_TOKEN}\" | docker login ghcr.io -u $ --password-stdin",
  env: {
    GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}",
  },
};

// Shared build steps list
local build_steps = [
  checkout_step,
  install_bazel('matrix'),
  // Native Tests
  bazel_step("linux_amd64", "matrix.host.platform_name == 'linux_amd64'"),
  bazel_step("linux_arm64", "matrix.host.platform_name == 'linux_arm64'"),
  // Cross Builds
  cross_build_step,
];

local lint_steps = [
  checkout_step,
  install_bazel('linux', 'arm64'),
  { name: "Reformat", run: "bazel run @com_github_buildbarn_bb_storage//tools:reformat" },
  { name: "Test style conformance", run: "git add . && git diff --exit-code HEAD --" },
  { name: "Golint", run: "bazel run @org_golang_x_lint//golint -- -set_exit_status $(pwd)/..." },
];

{
  "backend.yaml": {
    name: "Build and test backend",
    on: {
      push: { branches: ["main"] },
      pull_request: { branches: ["main"] },
    },
    jobs: {
      build_and_test: {
        name: "build_and_test ${{ matrix.host.platform_name }}",
        "runs-on": "${{ matrix.host.runs_on }}",
        strategy: { matrix: build_matrix },
        steps: build_steps,
      },
      lint: {
        name: "lint",
        "runs-on": "ubuntu-24.04-arm",
        steps: lint_steps,
      },
    },
  },

  "publish-docker.yaml": {
    name: "Build and publish docker images",
    on: {
      push: { branches: ["main"] },
      workflow_dispatch: null,
    },
    permissions: { contents: "read", "id-token": "write" },
    jobs: {
      publish: {
        name: "Publish bb-portal images",
        "runs-on": "ubuntu-24.04-arm",
        steps: [
          checkout_step,
          docker_credentials_step,
          install_bazel('linux', 'arm64'),
          {
            name: "Build and push backend",
            run: "bazel run --stamp //cmd/bb_portal:bb_portal_container_push"
          },
          {
            name: "Set tag variables",
            run: |||
              echo "TIMESTAMP=$(TZ=UTC date --date "@$(git show -s --format=%ct HEAD)" +%Y%m%dT%H%M%SZ)" >> $GITHUB_ENV
              echo "SHA_SHORT=$(git rev-parse --short HEAD)" >> $GITHUB_ENV
            |||
          },
          {
            name: "Build and push frontend",
            uses: "docker/build-push-action@v4",
            with: {
              context: "frontend",
              file: "./frontend/Dockerfile",
              push: true,
              tags: "ghcr.io/buildbarn/bb-portal-frontend:${{ env.TIMESTAMP }}-${{ env.SHA_SHORT }}",
            },
          },
        ],
      },
    },
  },

  "frontend.yaml": {
    name: "Build and test frontend",
    on: {
      push: { branches: ["main"] },
      pull_request: { branches: ["main"] },
    },
    jobs: {
      build_and_test: {
        name: "Build and test frontend",
        "runs-on": "ubuntu-24.04-arm",
        steps: [
          checkout_step,
          {
            name: "Build docker image",
            uses: "docker/build-push-action@v4",
            with: {
              context: "frontend/",
              file: "./frontend/Dockerfile",
              push: false,
              tags: "ghcr.io/buildbarn/bb-portal-frontend",
            },
          },
        ],
      },
    },
  },
}
