// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.

// This example config is made to be used with the docker-compose setup in
// [bb-deployments](https://github.com/buildbarn/bb-deployments), i.e. it
// assumes that the following services are running:
// - A Buildbarn scheduler, accessible at localhost:8984
// - A Buildbarn frontend, accessible at localhost:8980

{
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

  actionCacheProxy: {
    client: {
      address: 'localhost:8980',
    },
    allowedOrigins: ['http://localhost:8081'],
    httpServers: [{
      listenAddresses: [':9434'],
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

  contentAddressableStorageProxy: {
    client: {
      address: 'localhost:8980',
    },
    allowedOrigins: ['http://localhost:8081'],
    httpServers: [{
      listenAddresses: [':9435'],
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

  initialSizeClassCacheProxy: {
    client: {
      address: 'localhost:8980',
    },
    allowedOrigins: ['http://localhost:8081'],
    httpServers: [{
      listenAddresses: [':9436'],
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
  fileSystemAccessCacheProxy: {
    client: {
      address: 'localhost:8980',
    },
    allowedOrigins: ['http://localhost:8081'],
    httpServers: [{
      listenAddresses: [':9437'],
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
