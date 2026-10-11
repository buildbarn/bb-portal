import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerDrainsPage } from "@/components/pages/SchedulerDrains";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { SizeClassQueueNameSchema } from "@/utils/invocationName";

const DrainSearchSchema = z.object({
  sizeClassQueueName: SizeClassQueueNameSchema,
});

export type DrainSearchParams = z.infer<typeof DrainSearchSchema>;

export const Route = createFileRoute("/scheduler/drain")({
  component: RouteComponent,
  validateSearch: (search) => DrainSearchSchema.parse(search),
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Drains",
          _ctx.match.search.sizeClassQueueName.platformQueueName
            .instanceNamePrefix,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { sizeClassQueueName } = Route.useSearch();

  return <SchedulerDrainsPage sizeClassQueueName={sizeClassQueueName} />;
}
