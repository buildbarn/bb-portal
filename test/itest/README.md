# Local bb-portal and Buildbarn RBE stack

This target starts the frontend and backend as independent services, plus
PostgreSQL, Jaeger, and a complete local Buildbarn remote-execution path in
Docker Compose. Docker must be running. PostgreSQL and the volume initializer
are small images built from this repository by Bazel. The first start also
pulls the commit-stamped Buildbarn releases and digest-pinned Jaeger and Ubuntu
runner images, so it takes longer than later starts.

For development, install [iBazel](https://github.com/bazelbuild/bazel-watcher)
on your `PATH`, then start the stack so only the service whose inputs changed
is restarted:

```sh
ibazel run --config=enable_reload //test/itest:bb_portal
```

The frontend and backend are separate `rules_itest` services. Vite handles its
own browser hot module reload, while backend changes restart only the Go
service. PostgreSQL and the Buildbarn cluster share one Compose lifecycle and
stay running when none of their inputs changed. A one-shot run is also available:

```sh
bazel run --config=enable_reload //test/itest:bb_portal
```

The same service graph can be started, health-checked, and stopped as a test:

```sh
bazel test //test/itest:bb_portal_test
```

The `requires-network` tag lives on that test target, where it configures the
test sandbox. Compose keeps explicit service dependencies so startup order is
deterministic.

To avoid clashing with services already running on the host, every port
published by Docker Compose uses a `+10000` offset. For example, PostgreSQL is
published as `15432:5432`, while services inside the isolated Compose network
continue to use port `5432`. The local bb-portal listeners use the same offset
convention.

The stack uses these fixed loopback ports:

- Vite development server: <http://127.0.0.1:5173>
- Buildbarn scheduler administration: <http://127.0.0.1:17982>
- Buildbarn Remote Execution/CAS/AC: `grpc://127.0.0.1:18980`
- Buildbarn build queue state: `grpc://127.0.0.1:18984`
- Jaeger UI: <http://127.0.0.1:26686>
- OpenTelemetry OTLP: `grpc://127.0.0.1:14317` and <http://127.0.0.1:14318>
- PostgreSQL: `127.0.0.1:15432`
- bb-portal UI/API: <http://127.0.0.1:18081>
- bb-portal BES: `grpc://127.0.0.1:18082`

In a second terminal, execute the hermetic smoke action remotely and publish
its BEP to the local portal:

```sh
bazel build \
  --config=local_rbe \
  --config=bb_portal_itest \
  --noremote_accept_cached \
  //test/itest:rbe_smoke
```

`--config=local_rbe` supplies the endpoint, `hardlinking` instance name, and a
Linux/amd64 execution platform matching the worker. On Linux, the stack also
starts a FUSE-backed worker so bb-portal displays both Buildbarn
build-directory modes. Docker Desktop for macOS cannot provide the shared
bind-mount propagation that worker requires, so the FUSE Compose profile is
disabled there. This lets Bazel select
Linux-compatible execution tools, including host-configured bootstrap tools,
while leaving the build's target platform unchanged. Jaeger receives traces
from bb-portal and every Buildbarn runtime component. It stores them in its
embedded Badger database with a one-year retention period. The `jaeger-data`
named volume preserves the database across Jaeger restarts and normal Compose
shutdowns. Use the service selector in its UI to distinguish the portal,
storage frontend, storage, scheduler, worker, and runner. Repositories that
need additional remote toolchains can combine the same config with their
normal platform settings.

Stop iBazel with Ctrl-C, then remove the Compose containers while retaining
PostgreSQL, Jaeger, and Buildbarn's named volumes for the next run:

```sh
bazel run //test/itest:buildbarn_compose -- down
```

To also delete PostgreSQL data, Jaeger traces, the central CAS, and the action
cache, run:

```sh
bazel run //test/itest:buildbarn_compose -- down --volumes
```

The FUSE worker uses a bind-mounted workspace in Bazel's output tree so mount
propagation reaches its runner. `bazel clean` removes that cache after the
Compose stack has stopped.
