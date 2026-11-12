import { createChannel, createClient } from "nice-grpc-web";
import {
  type BuildQueueStateClient,
  BuildQueueStateDefinition,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

const clients = new Map<string, BuildQueueStateClient>();

export function getBuildQueueStateClient(
  schedulerId: string,
): BuildQueueStateClient {
  let client = clients.get(schedulerId);
  if (!client) {
    client = createClient(
      BuildQueueStateDefinition,
      createChannel(
        `/api/v1/schedulers/${encodeURIComponent(schedulerId)}/grpcweb`,
      ),
    );
    clients.set(schedulerId, client);
  }
  return client;
}
