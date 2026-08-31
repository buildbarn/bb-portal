# Local bb-portal and Buildbarn RBE stack

This target starts the frontend and backend as independent services, plus
PostgreSQL, Jaeger, and a complete local Buildbarn remote-execution path in
Docker Compose. Docker must be running. The first start pulls
the pinned Buildbarn and Ubuntu runner images, so it takes longer than later
starts.

For development, start the stack with iBazel so only the service whose inputs
changed is restarted:

```sh
~/.local/bin/ibazel run --config=enable_reload //test/itest:bb_portal
```

The frontend and backend are separate `rules_itest` services. Vite handles its
own browser hot module reload, while backend changes restart only the Go
service. PostgreSQL and the Buildbarn cluster share one Compose lifecycle and
stay running when none of their inputs changed. A one-shot run is also available:

```sh
bazel run --config=enable_reload //test/itest:bb_portal
```

The stack uses these fixed loopback ports:

- Vite development server: <http://127.0.0.1:5173>
- Buildbarn scheduler administration: <http://127.0.0.1:7982>
- Buildbarn Remote Execution/CAS/AC: `grpc://127.0.0.1:8980`
- Buildbarn build queue state: `grpc://127.0.0.1:8984`
- Jaeger UI: <http://127.0.0.1:16686>
- OpenTelemetry OTLP: `grpc://127.0.0.1:4317` and <http://127.0.0.1:4318>
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
Linux/amd64 execution platform matching the worker. This lets Bazel select
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

To also delete PostgreSQL data, Jaeger traces, the local CAS, and the action
cache, run:

```sh
bazel run //test/itest:buildbarn_compose -- down --volumes
```
