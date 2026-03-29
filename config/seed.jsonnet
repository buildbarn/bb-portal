// For documentation, see pkg/proto/configuration/bb_seed/bb_seed.proto.
//
// Usage: bb_seed config/seed.jsonnet
//
// NOTE: Stop bb_portal before seeding, as the cleanup job may delete
// seeded data that falls outside the invocation retention window.

{
  global: {
    tracing: {
      backends: [
        {
          otlpSpanExporter: {
            address: "localhost:4317"
          },
          simpleSpanProcessor: {}
        }
      ],
      resourceAttributes: [
        {
          key: "service.name",
          value: {
            stringValue: "bb-seed" 
          }
        }
      ],
      sampler: {
        always: {}
      }
    },
    diagnosticsHttpServer: {
      httpServers: [{
        listenAddresses: [':9981'],
        authenticationPolicy: { allow: {} },
      }],
      enablePrometheus: true,
      enablePprof: true,
      enableActiveSpans: true,
    },
  },
  database: {
    postgres: {
      connectionString: 'postgresql://app:password@localhost:5432/app',
    },
    connectionPoolConfiguration: {
      maxOpenConnections: 10,
      maxIdleConnections: 10,
      connectionMaxLifetime: '120s',
      connectionMaxIdleTime: '30s',
    },
  },
  instances: 10,
  users: 10,
  invocationsPerInstance: 100,
  targetsPerInvocation: 200,
  timeSpan: '432000s',
}
