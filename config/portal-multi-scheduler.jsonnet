// Minimal database-free Portal with two schedulers routed by instance name
// prefix, using the same map shape as bb_storage's `schedulers`.
// Start the frontend development server separately on port 5173.
// These are local example addresses, not a declaration of live inventory.
{
  global: {},
  maximumMessageSizeBytes: 2 * 1024 * 1024,
  httpServers: [{
    listenAddresses: [':8081'],
    authenticationPolicy: { allow: {} },
  }],
  schedulerServiceConfiguration: {
    schedulers: {
      // The scheduler sees "mars" for "ci/aws/mars" (bb_storage's default).
      'ci/aws': { endpoint: { address: 'localhost:8984' } },
      // The scheduler sees "ci/gcp/..." unchanged (perfect forwarding).
      'ci/gcp': {
        endpoint: { address: 'localhost:9984' },
        addInstanceNamePrefix: 'ci/gcp',
      },
    },
    // Authorizers always see Portal's instance names, e.g. "ci/aws/mars".
    readAuthorizer: { allow: {} },
    killOperationsAuthorizer: { deny: {} },
    listOperationsPageSize: 500,
  },
  frontendServiceConfiguration: {
    frontendSource: { proxy: 'http://localhost:5173' },
    frontendConfig: {
      companyName: 'Example Co',
      featureFlags: { scheduler: {} },
    },
  },
}
