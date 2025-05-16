// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.

// This example config is made to be used with the docker-compose setup in
// [bb-deployments](https://github.com/buildbarn/bb-deployments), i.e. it
// assumes that the following services are running:
// - A Buildbarn scheduler, accessible at localhost:8984
// - A Buildbarn frontend, accessible at localhost:8980

{
  frontendProxyUrl: 'http://localhost:3000',
  allowedOrigins: ['http://localhost:3000'],

  serveFilesCasConfiguration: {
    grpc: { address: 'localhost:8980' },
  },
  maximumMessageSizeBytes: 2 * 1024 * 1024,

  httpServers: [{
    listenAddresses: [':8081'],
    authenticationPolicy: {
      allow: {
        public: {
          user: 'FooBar',
        },
        private: {
          groups: ['admin'],
          instances: ['fuse', 'testingQueue'],
          email: 'foo@example.com',
        },
      },
    },
  }],
  grpcServers: [{
    listenAddresses: [':8082'],
    authenticationPolicy: { allow: {} },
    maximumReceivedMessageSizeBytes: 10 * 1024 * 1024,
  }],

  instanceNameAuthorizer: {
    jmespathExpression: |||
      contains(authenticationMetadata.private.instances, instanceName)
      || instanceName == ''
    |||,
  },

  buildQueueStateClient: {
    address: 'localhost:8984',
  },

  actionCacheClient: {
    address: 'localhost:8980',
  },

  contentAddressableStorageClient: {
    address: 'localhost:8980',
  },

  initialSizeClassCacheClient: {
    address: 'localhost:8980',
  },

  fileSystemAccessCacheClient: {
    address: 'localhost:8980',
  },

  listOperationsPageSize: 500,
}
