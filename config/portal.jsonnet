// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.

// This example config is made to be used with the docker-compose setup in
// [bb-deployments](https://github.com/buildbarn/bb-deployments), i.e. it
// assumes that the following services are running:
// - A Buildbarn scheduler, accessible at localhost:8984
// - A Buildbarn frontend, accessible at localhost:8980

// The application consists of 3 different services:
//  - BesService: A service that provides access to Buildbarn Execution Service
//    (BES) insights.
//  - BrowserService: A service that allows you to browse the contents of the
//    content addressable storage (CAS) and action cache.
//  - SchedulerService: A service that shows the state of the Buildbarn
//    scheduler
//
// Each service can be disabled by not setting the corresponding configuration.
// At least one service should be configured, otherwise the portal will not
// do anything useful.

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

  // The BesService can be disabled by not setting this field.
  besServiceConfiguration: {
    grpcServers: [{
      listenAddresses: [':8082'],
      authenticationPolicy: { allow: {} },
      maximumReceivedMessageSizeBytes: 10 * 1024 * 1024,
    }],
    database: {
      sqlite: {
<<<<<<< HEAD
        url: "file:buildportal.db?_journal=WAL&_fk=1"
=======
        connectionString: 'file:buildportal.db?_journal=WAL&_fk=1',
>>>>>>> b6926ba (fixup! Move configuration from CLI args to jsonnet-file)
      },
    },
    blobArchiveFolder: "./blob-archive/",
  },

  // The BrowserService can be disabled by not setting this field.
  browserServiceConfiguration: {
    contentAddressableStorage: { grpc: { address: 'localhost:8980' } },
    actionCache: { grpc: { address: 'localhost:8980' } },
    initialSizeClassCache: { grpc: { address: 'localhost:8980' } },
    fileSystemAccessCache: { grpc: { address: 'localhost:8980' } },
  },

  // The SchedulerService can be disabled by not setting this field.
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
