# Local bb-portal integration stack

Start a hermetic PostgreSQL server and bb-portal with its production frontend
embedded:

```sh
bazel run //test/itest:bb_portal
```

The service remains in the foreground. It uses these fixed loopback ports:

- PostgreSQL: `15432`
- Frontend/API: <http://127.0.0.1:18081>
- BES: `grpc://127.0.0.1:18082`

In a second terminal, send a Bazel invocation to it with the repository's
named Bazel configuration:

```sh
bazel build --config=bb_portal_itest //test/itest/postgres_healthcheck
```

Any `build`, `test`, or other command that supports the BES flags can use that
configuration. It also enables `--build_event_publish_all_actions` so successful
and cached actions are available in the invocation's Actions tab. Once the
upload finishes, open
<http://127.0.0.1:18081/bazel-invocations/>. Stop the stack with Ctrl-C; both
bb-portal and PostgreSQL are shut down by `rules_itest`, and the temporary
database is removed.
