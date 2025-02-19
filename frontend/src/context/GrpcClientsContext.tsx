import { BuildQueueStateClient } from '@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate';
import { createContext, useContext } from 'react';

interface GrpcClientsContextState {
  buildQueueStateClient: BuildQueueStateClient;
}

// biome-ignore lint/style/noNonNullAssertion: We want to throw an error if the context is used without provider, instead of failing silently.
export const GrpcClientsContext = createContext<GrpcClientsContextState>(null!);

export const useGrpcClients = () => useContext(GrpcClientsContext);
