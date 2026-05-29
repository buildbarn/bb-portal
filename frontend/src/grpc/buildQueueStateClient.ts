import { createChannel, createClient } from "nice-grpc-web";
import {
  type BuildQueueStateClient,
  BuildQueueStateDefinition,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

const grpcChannel = createChannel("/api/v1/grpcweb");

export const buildQueueStateClient: BuildQueueStateClient = createClient(
  BuildQueueStateDefinition,
  grpcChannel,
);
