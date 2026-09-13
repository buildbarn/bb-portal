local storage = { grpc: { client: { address: 'storage:8981' } } };

{
  blobstore: {
    contentAddressableStorage: storage,
    actionCache: {
      completenessChecking: {
        backend: storage,
        maximumTotalTreeSizeBytes: 64 * 1024 * 1024,
      },
    },
  },
  fileSystemAccessCache: storage,
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
        listenAddresses: [':9980'],
        authenticationPolicy: { allow: {} },
      }],
      enablePrometheus: true,
      enablePprof: true,
      enableActiveSpans: true,
    },
  },
}
