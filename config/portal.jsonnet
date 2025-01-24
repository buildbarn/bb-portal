// For documentation, see pkg/proto/configuration/bb_portal/bb_portal.proto.
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
}
