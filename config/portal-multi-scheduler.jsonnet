// Minimal database-free Portal with two independent BuildQueueState servers.
// Start the frontend development server separately on port 5173.
// These are local example addresses, not a declaration of live inventory.
local scheduler(id, label, address) = {
  id: id,
  label: label,
  configuration: {
    buildQueueStateClient: { address: address },
    readAuthorizer: { allow: {} },
    killOperationsAuthorizer: { deny: {} },
    listOperationsPageSize: 500,
  },
};
{
  global: {},
  maximumMessageSizeBytes: 2 * 1024 * 1024,
  httpServers: [{
    listenAddresses: [':8081'],
    authenticationPolicy: { allow: {} },
  }],
  namedSchedulers: [
    scheduler('ci', 'CI', 'localhost:8984'),
    scheduler('interactive', 'Interactive', 'localhost:9984'),
  ],
  defaultSchedulerId: 'ci',
  frontendServiceConfiguration: {
    frontendSource: { proxy: 'http://localhost:5173' },
    frontendConfig: {
      companyName: 'Example Co',
      featureFlags: { scheduler: {} },
    },
  },
}
