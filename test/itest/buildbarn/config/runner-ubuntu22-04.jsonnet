local common = import 'common.libsonnet';

{
  buildDirectoryPath: '/worker/build',
  global: common.global('bb-runner'),
  grpcServers: [{
    listenPaths: ['/worker/runner'],
    authenticationPolicy: { allow: {} },
  }],
}
