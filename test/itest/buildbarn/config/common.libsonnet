{
  blobstore: {
    contentAddressableStorage: {
      sharding: {
        shards: {
          '0': {
            backend: { grpc: { client: { address: 'storage-0:8981' } } },
            weight: 1,
          },
          '1': {
            backend: { grpc: { client: { address: 'storage-1:8981' } } },
            weight: 1,
          },
        },
      },
    },
    actionCache: {
      completenessChecking: {
        backend: {
          sharding: {
            shards: {
              '0': {
                backend: { grpc: { client: { address: 'storage-0:8981' } } },
                weight: 1,
              },
              '1': {
                backend: { grpc: { client: { address: 'storage-1:8981' } } },
                weight: 1,
              },
            },
          },
        },
        maximumTotalTreeSizeBytes: 64 * 1024 * 1024,
      },
    },
  },
  fileSystemAccessCache: {
    sharding: {
      shards: {
        '0': {
          backend: { grpc: { client: { address: 'storage-0:8981' } } },
          weight: 1,
        },
        '1': {
          backend: { grpc: { client: { address: 'storage-1:8981' } } },
          weight: 1,
        },
      },
    },
  },
  browserUrl: 'http://127.0.0.1:18081/browser',
  maximumMessageSizeBytes: 16 * 1024 * 1024,
  global(serviceName): {
    tracing: {
      backends: [{
        otlpSpanExporter: {
          address: 'jaeger:4317',
        },
        batchSpanProcessor: {},
      }],
      resourceAttributes: [{
        key: 'service.name',
        value: {
          stringValue: serviceName,
        },
      }],
      sampler: {
        always: {},
      },
    },
    diagnosticsHttpServer: {
      httpServers: [{
        listenAddresses: [':80'],
        authenticationPolicy: { allow: {} },
      }],
      enablePrometheus: true,
      enablePprof: true,
      enableActiveSpans: true,
    },
  },
}
