import { createFileRoute } from "@tanstack/react-router";
import { SchedulerPage } from "@/components/pages/Scheduler";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/scheduler/")({
  component: SchedulerPage,
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Scheduler"]) }] }),
});
