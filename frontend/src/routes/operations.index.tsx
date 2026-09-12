import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import {
  getExecutionStageFromOperationStatus,
  OperationStatus,
} from "@/components/OperationFilterSelector";
import { OperationsPage, PAGE_SIZE } from "@/components/pages/Operations";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { generatePageTitle } from "@/utils/generatePageTitle";

const OperationsFilterSchema = z
  .object({
    "@type": z.string(),
    toolInvocationId: z.uuid().optional(),
    correlatedInvocationsId: z.uuid().optional(),
  })
  .optional()
  .refine(
    (data) =>
      (data?.toolInvocationId && !data?.correlatedInvocationsId) ||
      (!data?.toolInvocationId && data?.correlatedInvocationsId),
    "Either toolInvocationId or correlatedInvocationsId must be provided, but not both",
  );

export type OperationsFilterParams = z.infer<typeof OperationsFilterSchema>;

const OperationsSearchSchema = z.object({
  filter: OperationsFilterSchema,
  statusFilter: z.enum(OperationStatus).optional(),
});

export const Route = createFileRoute("/operations/")({
  component: RouteComponent,
  validateSearch: (search) => OperationsSearchSchema.parse(search),
  loaderDeps: ({ search: { filter, statusFilter } }) => ({
    filter,
    statusFilter: statusFilter ?? OperationStatus.ALL,
  }),
  loader: async ({ deps }) => {
    const invocationFilter = deps.filter
      ? {
          typeUrl: deps.filter["@type"],
          value: RequestMetadata.encode(
            RequestMetadata.fromPartial(deps.filter),
          ).finish(),
        }
      : undefined;

    const response = await buildQueueStateClient.listOperations({
      pageSize: PAGE_SIZE,
      filterInvocationId: invocationFilter,
      filterStage: getExecutionStageFromOperationStatus(deps.statusFilter),
    });

    return { operations: response.operations };
  },
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Operations"]) }] }),
});

function RouteComponent() {
  const navigate = Route.useNavigate();
  const { filter, statusFilter } = Route.useLoaderDeps();
  const { operations } = Route.useLoaderData();

  const onStatusFilterChange = (value: OperationStatus): void => {
    navigate({
      from: Route.id,
      to: ".",
      search: (prev): typeof prev => ({
        ...prev,
        statusFilter: value,
      }),
    });
  };

  return (
    <OperationsPage
      filter={filter}
      statusFilter={statusFilter}
      onStatusFilterChange={onStatusFilterChange}
      operations={operations}
    />
  );
}
