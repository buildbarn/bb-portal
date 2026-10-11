import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerInvocationPage } from "@/components/pages/SchedulerInvocation";
import { QueuedOperationsCursorSchema } from "@/components/QueuedOperationsTable/utils";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { InvocationNameSchema } from "@/utils/invocationName";

export enum InvocationChildrenFilter {
  QUEUED = "queued",
  ACTIVE = "active",
  ALL = "all",
}

const InvocationSearchSchema = z.object({
  invocationName: InvocationNameSchema,
  childrenFilter: z
    .enum(InvocationChildrenFilter)
    .catch(InvocationChildrenFilter.ACTIVE),
  cursor: QueuedOperationsCursorSchema.optional(),
});

export type InvocationSearchParams = z.infer<typeof InvocationSearchSchema>;

export const Route = createFileRoute("/scheduler/invocation")({
  component: RouteComponent,
  validateSearch: (search) => InvocationSearchSchema.parse(search),
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          _ctx.match.search.invocationName.sizeClassQueueName.platformQueueName
            .instanceNamePrefix,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { invocationName, childrenFilter, cursor } = Route.useSearch();

  return (
    <SchedulerInvocationPage
      invocationName={invocationName}
      childrenFilter={childrenFilter}
      cursor={cursor}
    />
  );
}
