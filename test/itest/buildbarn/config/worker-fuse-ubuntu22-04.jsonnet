local common = import 'common.libsonnet';

{
  blobstore: {
    actionCache: common.blobstore.actionCache,
    contentAddressableStorage: {
      readCaching: {
        slow: common.blobstore.contentAddressableStorage,
        fast: {
          'local': {
            keyLocationMapOnBlockDevice: {
              file: {
                path: '/worker/cas/key_location_map',
                sizeBytes: 32 * 1024 * 1024,
              },
            },
            keyLocationMapMaximumGetAttempts: 16,
            keyLocationMapMaximumPutAttempts: 64,
            oldBlocks: 1,
            currentBlocks: 3,
            newBlocks: 1,
            blocksOnBlockDevice: {
              source: {
                file: {
                  path: '/worker/cas/blocks',
                  sizeBytes: 2 * 1024 * 1024 * 1024,
                },
              },
              spareBlocks: 1,
              dataIntegrityValidationCache: {
                cacheSize: 50000,
                cacheDuration: '14400s',
                cacheReplacementPolicy: 'LEAST_RECENTLY_USED',
              },
            },
            persistent: {
              stateDirectoryPath: '/worker/cas/persistent_state',
              minimumEpochInterval: '300s',
            },
          },
        },
        replicator: { deduplicating: { 'local': {} } },
      },
    },
  },
  browserUrl: common.browserUrl,
  maximumMessageSizeBytes: common.maximumMessageSizeBytes,
  scheduler: { address: 'scheduler:8983' },
  global: common.global('bb-worker-fuse'),
  buildDirectories: [{
    virtual: {
      maximumExecutionTimeoutCompensation: '3600s',
      shuffleDirectoryListings: true,
      maximumWritableFileUploadDelay: '60s',
      mount: {
        mountPath: '/worker/build',
        fuse: {
          directoryEntryValidity: '300s',
          inodeAttributeValidity: '300s',
          allowOther: true,
          mountMethod: 'DIRECT',
        },
      },
    },
    runners: [{
      endpoint: { address: 'unix:///worker/runner' },
      concurrency: 4,
      instanceNamePrefix: 'fuse',
      platform: {
        properties: [
          { name: 'OSFamily', value: 'linux' },
          { name: 'container-image', value: 'docker://ghcr.io/catthehacker/ubuntu:act-22.04@sha256:dd7654ffb01d5b7b54b23b9ce928a1f7f2d08c7b3d7e320b6574b55d7ccde78b' },
        ],
      },
      maximumFilePoolFileCount: 10000,
      maximumFilePoolSizeBytes: 256 * 1024 * 1024,
      workerId: { id: 'fuse' },
    }],
  }],
  filePool: {
    blockDevice: {
      file: {
        path: '/worker/filepool',
        sizeBytes: 1024 * 1024 * 1024,
      },
    },
  },
  inputDownloadConcurrency: 10,
  outputUploadConcurrency: 10,
  directoryCache: {
    maximumCount: 1000,
    maximumSizeBytes: 1024 * 1024,
    cacheReplacementPolicy: 'LEAST_RECENTLY_USED',
  },
  prefetching: {
    fileSystemAccessCache: common.fileSystemAccessCache,
    bloomFilterBitsPerPath: 14,
    bloomFilterMaximumSizeBytes: 65536,
  },
}
