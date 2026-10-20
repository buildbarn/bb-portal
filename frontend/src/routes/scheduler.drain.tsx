import { createFileRoute, redirect } from "@tanstack/react-router";
import { z } from "zod";
import { SchedulerDrainsPage } from "@/components/pages/SchedulerDrains";
import { generatePageTitle } from "@/utils/generatePageTitle";
import {
  SizeClassQueueNameSchema,
  type SizeClassQueueNameSearch,
} from "@/utils/sizeClassQueueName";

const DrainSearchSchema = z.object({
  sizeClassQueueName: SizeClassQueueNameSchema.optional().catch(undefined),
});

export type DrainSearchParams = {
  sizeClassQueueName: SizeClassQueueNameSearch;
};

export const Route = createFileRoute("/scheduler/drain")({
  component: RouteComponent,
  validateSearch: (search) => DrainSearchSchema.parse(search),
  beforeLoad: ({ search }) => {
    if (!search.sizeClassQueueName) {
      throw redirect({ to: "/scheduler" });
    }
  },
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

  // beforeLoad redirects without a queue.
  if (!sizeClassQueueName) {
    return null;
  }

  return <SchedulerDrainsPage sizeClassQueueName={sizeClassQueueName} />;
}
