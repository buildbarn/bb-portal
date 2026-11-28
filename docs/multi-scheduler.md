# Multiple schedulers

bb_storage can route Execute requests to several schedulers by instance
name prefix (`schedulers` in its configuration). Portal can show those
schedulers using the same map shape, so both can be generated from one
libsonnet file. See [the database-free example](../config/portal-multi-scheduler.jsonnet).

```jsonnet
schedulerServiceConfiguration: {
  schedulers: {
    'ci/aws': { endpoint: { address: 'scheduler-aws:8984' } },
    'ci/gcp': { endpoint: { address: 'scheduler-gcp:8984' }, addInstanceNamePrefix: 'ci/gcp' },
  },
  readAuthorizer: { ... },
  killOperationsAuthorizer: { ... },
  listOperationsPageSize: 500,
},
```

Each key is an instance name prefix, and the value is a
`buildbarn.configuration.builder.SchedulerConfiguration`, the message
bb_storage uses. Matching works the same way: the longest matching prefix
wins, the matched prefix is removed, and `addInstanceNamePrefix` is
prepended before the request is forwarded. Use the same keys and
`addInstanceNamePrefix` values as your frontends.

There is one difference: `endpoint` must be the scheduler's BuildQueueState
API (bb_scheduler's `buildQueueStateGrpcServers`), not the client API the
frontends use. A shared libsonnet helper can take the port as a parameter.

The single `buildQueueStateClient` form still works and is equivalent to a
`schedulers` map with one entry for the empty prefix. The two forms cannot be
combined.

## Authorization

`readAuthorizer` and `killOperationsAuthorizer` apply to every scheduler.
Portal rewrites instance names in responses back into its own namespace
before authorizing, so policies are written against the names clients use
(`ci/aws/mars`), even when the scheduler itself only sees `mars`.

If one scheduler is listed under several prefixes, each prefix only shows
the queues and operations within it. Portal never attributes a name outside
a prefix to that prefix.

## Web UI

The Scheduler and Operations pages show one scheduler at a time. With a
single scheduler there is no selector. With several, the page asks you to
choose a prefix, and the URL keeps the choice (`?scheduler=ci/aws`).

Requests that carry an instance name, such as listing the workers of a
queue, are routed by that name. Requests that do not, such as listing
queues or looking up an operation, are routed by the selected prefix. The
browser sends it in the `X-Bb-Portal-Scheduler-Instance-Name` header. API
clients that call the gRPC-Web endpoint directly must send it whenever more
than one scheduler is configured. If the header is missing or names an
unknown scheduler, `GetOperation` returns `NotFound` and the kill-permission
endpoint returns `{"allowed": false}`, consistent with authorization masking.
