# UnsignedLong manual validation environment

Run the complete environment with:

```sh
bazel run //test/manual/unsigned_long
```

The target builds and embeds the frontend, starts an ephemeral Postgres server,
generates a bb-portal configuration using the assigned ports, starts bb-portal,
uploads `unsigned_long_values.bep.ndjson`, and verifies the exact values through
GraphQL. The final startup log prints the invocation metrics URL to open in a
browser.

Press Ctrl-C to stop bb-portal and Postgres and remove their temporary data.
