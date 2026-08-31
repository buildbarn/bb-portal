local common = import 'common.libsonnet';

{
  blobstore: common.blobstore,
  browserUrl: common.browserUrl,
  maximumMessageSizeBytes: common.maximumMessageSizeBytes,
  scheduler: { address: 'scheduler:8983' },
  global: common.global('bb-worker'),
  buildDirectories: [{
    native: {
      buildDirectoryPath: '/worker/build',
      cacheDirectoryPath: '/worker/cache',
      maximumCacheFileCount: 10000,
      maximumCacheSizeBytes: 1024 * 1024 * 1024,
      cacheReplacementPolicy: 'LEAST_RECENTLY_USED',
    },
    runners: [{
      endpoint: { address: 'unix:///worker/runner' },
      concurrency: 4,
      instanceNamePrefix: 'hardlinking',
      platform: {
        properties: [
          { name: 'OSFamily', value: 'linux' },
          { name: 'container-image', value: 'docker://ghcr.io/catthehacker/ubuntu:act-22.04@sha256:dd7654ffb01d5b7b54b23b9ce928a1f7f2d08c7b3d7e320b6574b55d7ccde78b' },
        ],
      },
      workerId: {
        datacenter: 'local',
        rack: 'docker-compose',
        slot: '0',
        hostname: 'bb-portal-itest',
      },
    }],
  }],
  inputDownloadConcurrency: 10,
  outputUploadConcurrency: 10,
  directoryCache: {
    maximumCount: 1000,
    maximumSizeBytes: 1024 * 1024,
    cacheReplacementPolicy: 'LEAST_RECENTLY_USED',
  },
}
