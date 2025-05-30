// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.

// This example config is made to be used with the docker-compose setup in
// [bb-deployments](https://github.com/buildbarn/bb-deployments), i.e. it
// assumes that the following services are running:
// - A Buildbarn scheduler, accessible at localhost:8984
// - A Buildbarn frontend, accessible at localhost:8980

{
  frontendProxyUrl: 'http://localhost:3000',
  allowedOrigins: ['http://localhost:3000'],

  httpServers: [{
    listenAddresses: [':8081'],
    authenticationPolicy: { allow: {} },
  }],

  instanceNameAuthorizer: {
    allow: {},
  },

  maximumMessageSizeBytes: 2 * 1024 * 1024,

  besServiceConfiguration: {
    grpcServers: [{
      listenAddresses: [':8082'],
      authenticationPolicy: { allow: {} },
      maximumReceivedMessageSizeBytes: 10 * 1024 * 1024,
    }],
  },

  browserServiceConfiguration: {
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
    serveFilesCasConfiguration: {
      grpc: { address: 'localhost:8980' },
    },
  },

  schedulerServiceConfiguration: {
    buildQueueStateClient: {
      address: 'localhost:8984',
    },
    killOperationsAuthorizer: {
      allow: {},
    },
    listOperationsPageSize: 500,
  },

}
