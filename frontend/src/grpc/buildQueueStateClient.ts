import { createChannel, createClient, Metadata } from "nice-grpc-web";
import {
  type BuildQueueStateClient,
  BuildQueueStateDefinition,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

// Routes requests that carry no instance name (listing queues, looking up an
// operation) to the selected scheduler. Requests that carry one are routed by
// it on the server. Must match buildqueuestateproxy.SelectedInstanceNameMetadataKey.
export const SCHEDULER_SELECTION_HEADER = "x-bb-portal-scheduler-instance-name";

const grpcChannel = createChannel("/api/v1/grpcweb");
const clients = new Map<string, BuildQueueStateClient>();

export function getBuildQueueStateClient(
  instanceNamePrefix: string,
): BuildQueueStateClient {
  let client = clients.get(instanceNamePrefix);
  if (!client) {
    client = createClient(BuildQueueStateDefinition, grpcChannel, {
      "*": {
        metadata: Metadata({
          [SCHEDULER_SELECTION_HEADER]: instanceNamePrefix,
        }),
      },
    });
    clients.set(instanceNamePrefix, client);
  }
  return client;
}
