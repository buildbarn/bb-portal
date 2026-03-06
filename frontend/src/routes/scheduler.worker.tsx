import { generatePageTitle } from '@/utils/generatePageTitle';
import { SchedulerWorkersPage } from '@/components/pages/SchedulerWorkers';
import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';

export enum WorkerListStatus {
    ALL = 'all',
    EXECUTING = 'executing',
}

const PlatformProperySchema = z.object({
  name: z.string(),
  value: z.string(),
})

const WorkerSearchSchema = z.object({
  workerStatusFilter: z.enum(WorkerListStatus).catch(WorkerListStatus.ALL),
  sizeClassQueueName: z.object({
    platformQueueName: z.object({
      instanceNamePrefix: z.string(),
      platform: z.object({
        properties: z.array(PlatformProperySchema),
      })
    }),
    sizeClass: z.number().int().nonnegative().default(0).catch(0),
}),
  cursor: z.record(z.string(), z.string()).optional(),
})

export type WorkerSearchParams = z.infer<typeof WorkerSearchSchema>

export const Route = createFileRoute('/scheduler/worker')({
  component: RouteComponent,
  validateSearch: (search) => WorkerSearchSchema.parse(search),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Workers", _ctx.match.search.sizeClassQueueName.platformQueueName.instanceNamePrefix]) }] }),
})

function RouteComponent() {
  const { workerStatusFilter, sizeClassQueueName, cursor } = Route.useSearch();

  return <SchedulerWorkersPage
    workerStatusFilter={workerStatusFilter}
    sizeClassQueueName={sizeClassQueueName}
    cursor={cursor}
  />
}
