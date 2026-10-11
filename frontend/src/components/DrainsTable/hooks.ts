import {
  type QueryClient,
  queryOptions,
  useMutation,
  useQueries,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";
import { modifyDrainsPermissionUrl } from "./utils";

const modifyDrainsPermissionQuery = (instanceNamePrefix: string) =>
  queryOptions({
    queryKey: ["modifyDrainsButtonState", instanceNamePrefix],
    queryFn: async (): Promise<boolean> => {
      const response = await fetch(
        modifyDrainsPermissionUrl(instanceNamePrefix),
      );
      return (await response.json()).allowed;
    },
  });

export const useCanModifyDrains = (
  instanceNamePrefix: string,
  enabled = true,
): boolean | undefined =>
  useQuery({ ...modifyDrainsPermissionQuery(instanceNamePrefix), enabled })
    .data;

export const useCanModifyDrainsByPrefix = (
  instanceNamePrefixes: string[],
  enabled = true,
): Map<string, boolean | undefined> => {
  const unique = [...new Set(instanceNamePrefixes)];
  const results = useQueries({
    queries: unique.map((instanceNamePrefix) => ({
      ...modifyDrainsPermissionQuery(instanceNamePrefix),
      enabled,
    })),
  });
  return new Map(unique.map((prefix, i) => [prefix, results[i].data]));
};

export const invalidateDrainQueries = (queryClient: QueryClient) => {
  queryClient.invalidateQueries({ queryKey: ["listDrains"] });
  queryClient.invalidateQueries({ queryKey: ["listPlatformQueues"] });
  queryClient.invalidateQueries({ queryKey: ["listWorkers"] });
};

export const useAddDrain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      sizeClassQueueName,
      workerIdPattern,
    }: {
      sizeClassQueueName: SizeClassQueueNameSearch;
      workerIdPattern: Record<string, string>;
    }) =>
      buildQueueStateClient.addDrain({ sizeClassQueueName, workerIdPattern }),
    onSuccess: () => invalidateDrainQueries(queryClient),
  });
};
