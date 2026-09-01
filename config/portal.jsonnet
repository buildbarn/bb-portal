// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.

// This example config is made to be used with the docker-compose setup in
// [bb-deployments](https://github.com/buildbarn/bb-deployments), i.e. it
// assumes that the following services are running:
// - A Buildbarn scheduler, accessible at localhost:8984
// - A Buildbarn frontend, accessible at localhost:8980
//
// It also assumes a Postgresql database is running on localhost:5432 with
// username `app`, password `app`, and database `app`

local githubActionsExtractor = importstr 'gh-actions.jmespath';
local gitlabExtractor = importstr 'gitlab.jmespath';
local semaphoreExtractor = importstr 'semaphore.jmespath';

local readAuthorizer = { allow: {} };
local writeAuthorizer = { allow: {} };
local adminAuthorizer = { allow: {} };

{
  global: {
    tracing: {
      backends: [
        {
          otlpSpanExporter: {
            address: 'localhost:4317',
          },
          batchSpanProcessor: {},
        },
      ],
      resourceAttributes: [
        {
          key: 'service.name',
          value: {
            stringValue: 'bb-portal',
          },
        },
      ],
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
  maximumMessageSizeBytes: 2 * 1024 * 1024,

  httpServers: [{
    listenAddresses: [':8081'],
    authenticationPolicy: { allow: {} },
  }],

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
    cleanupConfiguration: {
      cleanupInterval: '60s',
      invocationMessageTimeout: '3600s',
      invocationRetention: '604800s',
      completedActionRetention: '604800s',
    },
  },

  contentAddressableStorage: {
    backend: { grpc: { client: { address: 'localhost:8980' } } },
    readAuthorizer: readAuthorizer,
  },
  actionCache: {
    backend: { grpc: { client: { address: 'localhost:8980' } } },
    readAuthorizer: readAuthorizer,
  },
  initialSizeClassCache: {
    backend: { grpc: { client: { address: 'localhost:8980' } } },
    readAuthorizer: readAuthorizer,
  },
  fileSystemAccessCache: {
    backend: { grpc: { client: { address: 'localhost:8980' } } },
    readAuthorizer: readAuthorizer,
  },

  blobstoreServiceConfiguration: {},

  besServiceConfiguration: {
    grpcServers: [{
      listenAddresses: [':8082'],
      authenticationPolicy: { allow: {} },
      maximumReceivedMessageSizeBytes: 10 * 1024 * 1024,
    }],
    publishAuthorizer: writeAuthorizer,
    enableBepFileUpload: true,
    saveDataLevel: { basicAndTarget: {} },
    minEventBatchDuration: '0.1s',
    invocationMetadataExtractor: {
      expression: githubActionsExtractor,
    },
    buildKey: 'build_id',
  },

  schedulerServiceConfiguration: {
    buildQueueStateClient: {
      address: 'localhost:8984',
    },
    readAuthorizer: readAuthorizer,
    killOperationsAuthorizer: adminAuthorizer,
    listOperationsPageSize: 500,
  },

  graphqlApiServiceConfiguration: {
    readAuthorizer: readAuthorizer,
  },

  frontendServiceConfiguration: {
    frontendSource: {
      // NOTE: In production, you should use the `embedded` option instead of
      // `proxy`.
      proxy: 'http://localhost:5173',
      // embedded: {},
    },
    frontendConfig: {
      companyName: 'Example Co',
      grpcBackendUrl: 'grpc://localhost:8082',
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
      footerContent: [
        {
          text: 'Buildteam',
          href: 'https://buildteamworld.slack.com/archives/CD6HZC750',
          icon: { slack: {} },
        },
      ],
      additionalBuildColumns: [
        { title: 'Repo', value_key: 'repo', url_key: 'repo_url' },
        { title: 'PR', value_key: 'pull_request', url_key: 'pull_request_url' },
        { title: 'Workflow', value_key: 'workflow', url_key: 'workflow_url' },
      ],
      additionalBuildInvocationColumns: [
        { title: 'Job', value_key: 'job' },
        { title: 'Action', value_key: 'action' },
      ],
    },
  },
}
