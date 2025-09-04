# Adding a new schema

Important: When adding a new schema, make sure to consider authorization for
the new schema.

The existing schemas are authorized by using the `instanceNameAuthorizer`. The
models `BazelInvocation` and `Blobs` contains the filed `instanceName` and the
rest of the schemas get their instance name from either `BazelInvocation` or
`Blobs`. This is implemented in `internal/graphql/auth/auth.go`.

If your new schema can be connected to `BazelInvocation` or `Blobs`, you can
add it similarly to the existing schemas. If not, consider if it would be
possible to add the instance name to the new schema, or if it does not need
authorization at all.
