import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerDrainsPage } from "@/components/pages/SchedulerDrains";
import { SchedulerQueueChooserPage } from "@/components/pages/SchedulerQueueChooser";
import { generatePageTitle } from "@/utils/generatePageTitle";
import {
  SizeClassQueueNameSchema,
  type SizeClassQueueNameSearch,
} from "@/utils/invocationName";

const DrainSearchSchema = z.object({
  sizeClassQueueName: SizeClassQueueNameSchema.optional().catch(undefined),
});

export type DrainSearchParams = {
  sizeClassQueueName: SizeClassQueueNameSearch;
};

export const Route = createFileRoute("/scheduler/drain")({
  component: RouteComponent,
  validateSearch: (search) => DrainSearchSchema.parse(search),
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Drains",
          _ctx.match.search.sizeClassQueueName?.platformQueueName
            .instanceNamePrefix,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { sizeClassQueueName } = Route.useSearch();

  if (!sizeClassQueueName) {
    return <SchedulerQueueChooserPage title="Drains" />;
  }

  return <SchedulerDrainsPage sizeClassQueueName={sizeClassQueueName} />;
}
