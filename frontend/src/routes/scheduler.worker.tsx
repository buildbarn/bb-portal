import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerWorkersPage } from "@/components/pages/SchedulerWorkers";
import { generatePageTitle } from "@/utils/generatePageTitle";
import {
  InvocationIdSchema,
  SizeClassQueueNameSchema,
} from "@/utils/invocationName";

export enum WorkerListStatus {
  ALL = "all",
  EXECUTING = "executing",
  IDLE_SYNCHRONIZING = "idleSynchronizing",
}

const WorkerSearchSchema = z.object({
  workerStatusFilter: z.enum(WorkerListStatus).catch(WorkerListStatus.ALL),
  sizeClassQueueName: SizeClassQueueNameSchema,
  invocationIds: z.array(InvocationIdSchema).optional(),
  cursor: z.record(z.string(), z.string()).optional(),
});

export type WorkerSearchParams = z.infer<typeof WorkerSearchSchema>;

export const Route = createFileRoute("/scheduler/worker")({
  component: RouteComponent,
  validateSearch: (search) => WorkerSearchSchema.parse(search),
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Workers",
          _ctx.match.search.sizeClassQueueName.platformQueueName
            .instanceNamePrefix,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { workerStatusFilter, sizeClassQueueName, invocationIds, cursor } =
    Route.useSearch();

  return (
    <SchedulerWorkersPage
      workerStatusFilter={workerStatusFilter}
      sizeClassQueueName={sizeClassQueueName}
      invocationIds={invocationIds}
      cursor={cursor}
    />
  );
}
