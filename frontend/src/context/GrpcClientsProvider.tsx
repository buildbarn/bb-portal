import {
  BuildQueueStateClient,
  BuildQueueStateDefinition,
} from '@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate';
import { env } from 'next-runtime-env';
import { createChannel, createClient } from 'nice-grpc-web';
import { ReactNode } from 'react';
import { GrpcClientsContext } from './GrpcClientsContext';

export interface GrpcClientsProviderProps {
  children: ReactNode;
}

const GrpcClientsProvider = ({ children }: GrpcClientsProviderProps) => {
  const buildQueueStateClient: BuildQueueStateClient = createClient(
    BuildQueueStateDefinition,
    createChannel(env("NEXT_PUBLIC_BB_BUILDQUEUESTATE_GRPC_BACKEND_URL") || ""),
  );

  return (
    <GrpcClientsContext.Provider
      value={{
        buildQueueStateClient,
      }}
    >
      {children}
    </GrpcClientsContext.Provider>
  );
};

export default GrpcClientsProvider;
