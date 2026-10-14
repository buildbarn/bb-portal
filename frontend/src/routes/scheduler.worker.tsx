import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerQueueChooserPage } from "@/components/pages/SchedulerQueueChooser";
import { SchedulerWorkersPage } from "@/components/pages/SchedulerWorkers";
import { generatePageTitle } from "@/utils/generatePageTitle";
import {
  InvocationIdSchema,
  SizeClassQueueNameSchema,
  type SizeClassQueueNameSearch,
} from "@/utils/invocationName";

export enum WorkerListStatus {
  ALL = "all",
  EXECUTING = "executing",
  IDLE_SYNCHRONIZING = "idleSynchronizing",
}

const WorkerSearchSchema = z.object({
  workerStatusFilter: z.enum(WorkerListStatus).catch(WorkerListStatus.ALL),
  sizeClassQueueName: SizeClassQueueNameSchema.optional().catch(undefined),
  invocationIds: z.array(InvocationIdSchema).optional(),
  cursor: z.record(z.string(), z.string()).optional(),
});

export type WorkerSearchParams = Omit<
  z.infer<typeof WorkerSearchSchema>,
  "sizeClassQueueName"
> & { sizeClassQueueName: SizeClassQueueNameSearch };

export const Route = createFileRoute("/scheduler/worker")({
  component: RouteComponent,
  validateSearch: (search) => WorkerSearchSchema.parse(search),
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Workers",
          _ctx.match.search.sizeClassQueueName?.platformQueueName
            .instanceNamePrefix,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { workerStatusFilter, sizeClassQueueName, invocationIds, cursor } =
    Route.useSearch();

  if (!sizeClassQueueName) {
    return <SchedulerQueueChooserPage title="Workers" />;
  }

  return (
    <SchedulerWorkersPage
      workerStatusFilter={workerStatusFilter}
      sizeClassQueueName={sizeClassQueueName}
      invocationIds={invocationIds}
      cursor={cursor}
    />
  );
}
