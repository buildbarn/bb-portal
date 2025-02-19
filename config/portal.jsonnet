// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.

// This example config is made to be used with the docker-compose setup in
// [bb-deployments](https://github.com/buildbarn/bb-deployments), i.e. it
// assumes that the following services are running:
// - A Buildbarn scheduler, accessible at localhost:8984

{
  httpServers: [{
    listenAddresses: [':8081'],
    authenticationPolicy: { allow: {} },
  }],
  grpcServers: [{
    listenAddresses: [':8082'],
    authenticationPolicy: { allow: {} },
    maximumReceivedMessageSizeBytes: 10*1024*1024
  }],

  instanceNameAuthorizer: {
    jmespathExpression: |||
      contains(authenticationMetadata.private.instances, instanceName)
      || instanceName == ''
    |||,
  },

  buildQueueStateProxy: {
    client: {
      address: 'localhost:8984',
    },
    allowedOrigins: ['http://localhost:8081'],
    httpServers: [{
      listenAddresses: [':9433'],
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
  },
}
