# Multiple scheduler routing integration tests

Run from the Portal root:

```sh
go test -race -v ./test/multischeduler
```

The tests start TCP gRPC BuildQueueState servers keyed by instance name
prefix, with colliding operation IDs and worker IDs. One scheduler sees
stripped instance names (bb_storage's default) and the other uses perfect
forwarding. The tests instantiate Portal's production scheduler service and
send binary HTTP/gRPC-Web requests through the actual service path, parsing
protobuf frames and status trailers. No database, cluster, external scheduler,
or image is needed.

Coverage includes routing by the instance name in a request, routing by the
web UI's selection header for requests without one, instance names rewritten
back into Portal's namespace, a single kill authorizer written against those
names, pagination, the empty prefix as a selection, single-scheduler
configurations without a selection, an unreachable scheduler at startup,
failed upstream isolation, and HTTP request cancellation propagating to the
upstream context. A separate test simulates the outer HTTP identity
authentication layer with a test-only trusted identity header and verifies
authentication context reaches both gRPC-Web reads and permission checks; it
does not introduce a production header authentication policy.

These tests do not launch a browser or validate a live deployment. The test
client factory uses local plaintext gRPC connections; TLS configuration is not
covered by this harness.
