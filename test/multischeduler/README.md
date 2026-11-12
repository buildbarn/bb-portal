# Independent scheduler routing integration tests

Run from the Portal root:

```sh
go test -race -v ./test/multischeduler
```

The tests start two independent TCP gRPC BuildQueueState servers with different
operation data but colliding operation IDs, worker IDs and queue names. They
instantiate Portal's production scheduler registry/router and send binary
HTTP/gRPC-Web requests through the actual service path, parsing protobuf frames
and status trailers. No database, cluster, external scheduler, or image is needed.

Coverage includes per-scheduler operations, workers, queues, pagination,
read authorization, cancellation permission/mutation agreement, legacy default
aliases, legacy and single-named configuration with automatic defaults, unknown
IDs, an unreachable default backend at startup, failed upstream isolation, and
HTTP request cancellation propagating to the upstream context. A separate test
simulates the outer HTTP identity authentication layer with a test-only trusted
identity header and
verifies authentication context reaches both scoped gRPC-Web reads and permission
checks; it does not introduce a production header authentication policy.

These tests do not launch a browser or validate a live deployment. The test
client factory uses local plaintext gRPC connections; TLS configuration is not
covered by this harness.
