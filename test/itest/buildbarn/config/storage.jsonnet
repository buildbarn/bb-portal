local common = import 'common.libsonnet';

local localBlobstore(path, keyLocationMapSizeBytes, blocksSizeBytes) = {
  'local': {
    keyLocationMapOnBlockDevice: {
      file: {
        path: path + '/key_location_map',
        sizeBytes: keyLocationMapSizeBytes,
      },
    },
    keyLocationMapMaximumGetAttempts: 16,
    keyLocationMapMaximumPutAttempts: 64,
    oldBlocks: 8,
    currentBlocks: 24,
    // AC and FSAC objects may be updated in place, which requires a single
    // writable generation. This setting is also valid for the CAS.
    newBlocks: 1,
    blocksOnBlockDevice: {
      source: {
        file: {
          path: path + '/blocks',
          sizeBytes: blocksSizeBytes,
        },
      },
      spareBlocks: 3,
    },
    persistent: {
      stateDirectoryPath: path + '/persistent_state',
      minimumEpochInterval: '300s',
    },
  },
};

{
  grpcServers: [{
    listenAddresses: [':8981'],
    authenticationPolicy: { allow: {} },
  }],
  maximumMessageSizeBytes: common.maximumMessageSizeBytes,
  global: common.global('bb-storage'),
  contentAddressableStorage: {
    // Prebuilt remote tools such as Node exceed the per-block limit of a
    // 4-GiB device with this block layout.
    backend: localBlobstore('/storage-cas', 32 * 1024 * 1024, 8 * 1024 * 1024 * 1024),
    getAuthorizer: { allow: {} },
    putAuthorizer: { allow: {} },
    findMissingAuthorizer: { allow: {} },
  },
  actionCache: {
    backend: localBlobstore('/storage-ac', 1024 * 1024, 20 * 1024 * 1024),
    getAuthorizer: { allow: {} },
    putAuthorizer: { allow: {} },
  },
  fileSystemAccessCache: {
    backend: localBlobstore('/storage-fsac', 1024 * 1024, 20 * 1024 * 1024),
    getAuthorizer: { allow: {} },
    putAuthorizer: { allow: {} },
  },
}
