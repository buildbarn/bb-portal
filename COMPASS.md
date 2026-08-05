# Compass Fork — bb-portal

This is Compass's fork of [buildbarn/bb-portal](https://github.com/buildbarn/bb-portal).
bb-portal is the BuildBarn component Compass plans to extensively customize, so the fork and
image pipeline are established early (TECH-23509) before customization begins.

---

## Repository Decisions

**Default branch is `compass/main`, not `main`.**
Keeps Compass work clearly separated from upstream. All Compass changes are merged to `compass/main`.
Upstream `main` is left untouched to simplify future upstream syncs.

**Branch protection on `compass/main`.**
Requires 1 approving review. Force pushes disabled. Prevents accidental overwrites of the default branch.

**All upstream GitHub Actions workflows are disabled.**
Compass uses CircleCI only. The upstream `.github/workflows/` files (`backend.yaml`, `frontend.yaml`,
`publish-docker.yaml`) were removed to avoid confusion and unwanted runs.

**No internal values hardcoded in the repo.**
This is a public fork. AWS account IDs, ECR endpoints, and role ARNs are never committed.
All Compass-specific values are injected at CI runtime via CircleCI project environment variables.

---

## Image Pipeline Decisions

**Bazel-based image build, not a Dockerfile.**
Uses `multiarch_go_image` + `image_push` from `rules_img` — the same toolchain the upstream
project uses. Keeps the build hermetic and consistent with how upstream publishes images.

**ECR registry and repo injected via Bazel workspace status stamping.**
`tools/compass-workspace-status.sh` reads `$ECR_REGISTRY` and `$ECR_REPO` from environment and
outputs them as `STABLE_COMPASS_ECR_REGISTRY` / `STABLE_COMPASS_ECR_REPO`. The `image_push`
target in `cmd/bb_portal/BUILD.bazel` consumes these via Go template syntax (`{{.STABLE_*}}`).
Nothing is hardcoded in the Bazel files.

**Multi-arch image (amd64 + arm64).**
`multiarch_go_image` builds for both architectures automatically. arm64 is included for future
Graviton (EKS) targeting. The manifest index is what gets tagged `:latest`.

**ECR push uses `sts:AssumeRole` on `ECRPowerUserRole`.**
Mirrors the pattern used across the urbancompass repo (`build_infra_img_ecr.sh`). The base IAM
credentials (project-level env vars in CircleCI) only need `sts:AssumeRole` permission; the
actual ECR push permission comes from the assumed role. Credentials are never printed to stdout —
`jq` output is redirected directly to `$BASH_ENV`.

**Role ARN is derived from `ECR_REGISTRY`, not a separate env var.**
The account ID is already embedded in the `ECR_REGISTRY` URL
(e.g. `149465543054.dkr.ecr.us-east-1.amazonaws.com`). The login script extracts it with
`cut -d. -f1` and constructs the role ARN inline — one less env var to manage, no internal IDs
in the repo.

**CI pipeline tags with both `:latest` and the commit SHA.**
`bazel run` pushes `:latest`. A subsequent `crane tag` adds the `$CIRCLE_SHA1` tag.
This makes it easy to pin to a specific build while always having a stable `:latest`.

---

## CI Environment Variables

These must be set as project-level environment variables in the CircleCI project
(not as a context):

| Variable | Description |
|---|---|
| `AWS_ACCESS_KEY_ID` | IAM user with `sts:AssumeRole` on `ECRPowerUserRole` — same creds used in the urbancompass repo's `aws_auth` |
| `AWS_SECRET_ACCESS_KEY` | As above |
| `AWS_DEFAULT_REGION` | `us-east-1` |
| `ECR_REGISTRY` | `<account-id>.dkr.ecr.<region>.amazonaws.com` — account ID is also used to derive the `ECRPowerUserRole` ARN |
| `ECR_REPO` | ECR repository name (currently `buildtools_bbportal`) |

---

## Syncing Upstream Changes

To pull in upstream changes from `buildbarn/bb-portal`:

```bash
git fetch upstream
git checkout compass/main
git merge upstream/main
# resolve any conflicts, then open a PR to compass/main
```

The upstream remote should point to `https://github.com/buildbarn/bb-portal`.
Compass-specific files (`COMPASS.md`, `.circleci/`, `tools/compass-*.sh`) will never conflict
with upstream as they don't exist there.
