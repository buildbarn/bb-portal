import { createFileRoute, redirect } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerWorkersPage } from "@/components/pages/SchedulerWorkers";
import { generatePageTitle } from "@/utils/generatePageTitle";

export enum WorkerListStatus {
  ALL = "all",
  EXECUTING = "executing",
}

const PlatformProperySchema = z.object({
  name: z.string(),
  value: z.string(),
});

const SizeClassQueueNameSchema = z.object({
  platformQueueName: z.object({
    instanceNamePrefix: z.string(),
    platform: z.object({
      properties: z.array(PlatformProperySchema),
    }),
  }),
  sizeClass: z.number().int().nonnegative().default(0).catch(0),
});

const WorkerSearchSchema = z.object({
  workerStatusFilter: z.enum(WorkerListStatus).catch(WorkerListStatus.ALL),
  sizeClassQueueName: SizeClassQueueNameSchema.optional().catch(undefined),
  cursor: z.record(z.string(), z.string()).optional(),
});

export type WorkerSearchParams = Omit<
  z.infer<typeof WorkerSearchSchema>,
  "sizeClassQueueName"
> & { sizeClassQueueName: z.infer<typeof SizeClassQueueNameSchema> };

export const Route = createFileRoute("/scheduler/worker")({
  component: RouteComponent,
  validateSearch: (search) => WorkerSearchSchema.parse(search),
  beforeLoad: ({ search }) => {
    if (!search.sizeClassQueueName) {
      throw redirect({ to: "/scheduler" });
    }
  },
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
  const { workerStatusFilter, sizeClassQueueName, cursor } = Route.useSearch();

  // beforeLoad redirects without a queue.
  if (!sizeClassQueueName) {
    return null;
  }

  return (
    <SchedulerWorkersPage
      workerStatusFilter={workerStatusFilter}
      sizeClassQueueName={sizeClassQueueName}
      cursor={cursor}
    />
  );
}
