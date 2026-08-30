// Standalone configuration for //test/itest:bb_portal.
//
// The fixed loopback ports intentionally make it easy for a separate Bazel
// invocation to stream its BEP to this service:
//   PostgreSQL: 127.0.0.1:15432
//   Vite:       http://127.0.0.1:5173
//   Buildbarn:  grpc://127.0.0.1:8980
//   Jaeger:     http://127.0.0.1:16686
//   Frontend:   http://127.0.0.1:18081
//   BES:        grpc://127.0.0.1:18082
{
  global: {
    tracing: {
      backends: [{
        otlpSpanExporter: {
          address: '127.0.0.1:4317',
        },
        batchSpanProcessor: {},
      }],
      resourceAttributes: [{
        key: 'service.name',
        value: {
          stringValue: 'bb-portal',
        },
      }],
      sampler: {
        always: {},
      },
    },
  },

  httpServers: [{
    listenAddresses: ['127.0.0.1:18081'],
    authenticationPolicy: { allow: {} },
  }],

  instanceNameAuthorizer: { allow: {} },
  maximumMessageSizeBytes: 16 * 1024 * 1024,

  contentAddressableStorage: { grpc: { client: { address: '127.0.0.1:8980' } } },
  actionCache: { grpc: { client: { address: '127.0.0.1:8980' } } },
  fileSystemAccessCache: { grpc: { client: { address: '127.0.0.1:8980' } } },

  besServiceConfiguration: {
    grpcServers: [{
      listenAddresses: ['127.0.0.1:18082'],
      authenticationPolicy: { allow: {} },
      maximumReceivedMessageSizeBytes: 64 * 1024 * 1024,
    }],
    database: {
      postgres: {
        connectionString: 'postgresql://app:password@127.0.0.1:15432/postgres?sslmode=disable',
      },
      connectionPoolConfiguration: {
        maxOpenConnections: 10,
        maxIdleConnections: 10,
        connectionMaxLifetime: '120s',
        connectionMaxIdleTime: '30s',
      },
    },
    enableBepFileUpload: true,
    enableGraphqlPlayground: true,
    saveDataLevel: { basicAndTarget: {} },
    databaseCleanupConfiguration: {
      cleanupInterval: '60s',
      invocationMessageTimeout: '3600s',
      invocationRetention: '86400s',
    },
    minEventBatchDuration: '0s',
  },

  schedulerServiceConfiguration: {
    buildQueueStateClient: {
      address: '127.0.0.1:8984',
    },
    killOperationsAuthorizer: { allow: {} },
    listOperationsPageSize: 500,
  },

  frontendServiceConfiguration: {
    frontendSource: { proxy: 'http://127.0.0.1:5173' },
    frontendConfig: {
      companyName: 'bb-portal integration test',
      grpcBackendUrl: 'grpc://127.0.0.1:18082',
      featureFlags: {
        home: {
          fileUpload: {},
          instructions: {},
        },
        bes: {
          pageBuilds: {},
          pageInvocations: {},
          pageTargets: {},
          pageTests: {},
          pageTrends: {},
        },
        browser: {},
        scheduler: {},
      },
    },
  },
}
