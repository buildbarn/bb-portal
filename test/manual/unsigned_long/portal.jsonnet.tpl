{
  global: {},

  httpServers: [{
    listenAddresses: ['127.0.0.1:@@HTTP_PORT@@'],
    authenticationPolicy: { allow: {} },
  }],

  instanceNameAuthorizer: { allow: {} },

  besServiceConfiguration: {
    database: {
      postgres: {
        connectionString: 'postgresql://postgres:postgres@127.0.0.1:@@POSTGRES_PORT@@/postgres?sslmode=disable',
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
      cleanupInterval: '86400s',
      invocationMessageTimeout: '3600s',
      invocationRetention: '3153600000s',
    },
  },

  frontendServiceConfiguration: {
    frontendSource: { embedded: {} },
    frontendConfig: {
      companyName: 'UnsignedLong validation',
      featureFlags: {
        home: { fileUpload: {} },
        bes: { pageInvocations: {} },
      },
    },
  },
}
