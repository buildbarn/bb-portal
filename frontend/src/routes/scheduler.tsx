import { createFileRoute } from "@tanstack/react-router";
import { SchedulerError, SchedulerLayout } from "@/components/SchedulerLayout";
import { schedulerContext } from "@/schedulers/context";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";

export const Route = createFileRoute("/scheduler")({
  component: SchedulerLayout,
  errorComponent: SchedulerError,
  beforeLoad: async ({ search }) => {
    requireFeature(env.featureFlags?.scheduler)();
    return schedulerContext(search.scheduler);
  },
});
