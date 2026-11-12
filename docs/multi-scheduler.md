# Named schedulers

Portal can serve several independent BuildQueueState APIs while displaying one
scheduler at a time. Scheduler IDs are Portal identities; instance names still
identify execution/storage namespaces. No database or scheduler protocol changes
are required.

See [the database-free example](../config/portal-multi-scheduler.jsonnet).
Each `namedSchedulers` entry has an `id`, a display `label`, and a `configuration`
containing the complete existing scheduler configuration: gRPC client settings,
read authorizer, cancellation authorizer, and pagination size. Keep existing TLS,
credentials and authorizers in that configuration when migrating.

IDs must match `[a-z0-9][a-z0-9_-]*` and contain at most 63 characters. They must be
unique. Labels default to the ID when omitted. A single scheduler is automatically
the default; multiple schedulers require `defaultSchedulerId`. Missing clients,
unknown defaults and mixed legacy/named configuration are startup errors. Omit
both scheduler configuration forms to disable scheduler APIs, and omit the
frontend scheduler feature flag to hide those pages.

The browser obtains IDs, labels and the default from `GET /api/v1/schedulers`.
Backend addresses, credentials and authorizers are not returned. The catalog is
loaded once per browser page; reload Portal after changing scheduler configuration.
A URL such as
`/scheduler?scheduler=interactive` keeps the selection across refreshes and shared
links. Tabs maintain their own selection through their URLs. Old URLs without a
selection resolve to the configured default. An explicit unknown ID produces an
error with a way to choose a configured scheduler.

Scheduler RPCs use
`/api/v1/schedulers/{id}/grpcweb/buildbarn.buildqueuestate.BuildQueueState/{method}`.
Cancellation permission checks use
`/api/v1/schedulers/{id}/checkPermissions/killOperation/{operationName}`.
The legacy RPC and permission endpoints remain aliases for the default scheduler.
Only configured IDs are accepted; the browser cannot supply a backend address.
Existing read and cancellation authorizers apply separately to each scheduler.

Changing the selector returns to the scheduler or operations list and resets
scheduler-specific filters and details. It never opens a matching operation ID on
another scheduler automatically. Action, log and output links retain their original
instance names and existing storage routes. Unavailable schedulers fail their own
requests without changing selection or preventing access to healthy schedulers.

## Migration and rollback

1. Build and release a Portal binary containing named-scheduler support before
   emitting the new configuration fields. Obtain the real image reference from
   CI; do not infer a tag from a branch name.
2. Replace `schedulerServiceConfiguration` with a `namedSchedulers` entry whose
   `configuration` copies the complete legacy object. Choose a stable ID and
   default. Do not supply both forms together.
3. Add independently verified BuildQueueState API addresses and authorizers.
   Use each scheduler's configured BuildQueueState API endpoint, which is
   separate from its execution API and admin UI. Verify that each ID points to
   the intended independent scheduler.
4. Validate each scheduler and scoped permission route before using a deployment.
   Preserve the existing cancellation and read authorization policies.

To roll back configuration, restore the complete legacy scheduler object and
remove `namedSchedulers` and `defaultSchedulerId` first. The new binary supports
that configuration. Only then roll back the binary to one that cannot parse the
new fields. Links with explicit IDs for removed schedulers will show an error;
legacy links continue to resolve to the configured default. Keep TLS/proxy and
storage configuration unchanged throughout.

An aggregate view, cross-environment federation and new storage routing are out
of scope.
