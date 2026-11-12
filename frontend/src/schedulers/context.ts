import { useRouteContext } from "@tanstack/react-router";
import { getBuildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { loadSchedulerRegistry, resolveScheduler } from "./registry";

export async function schedulerContext(requestedId?: unknown) {
  const schedulerRegistry = await loadSchedulerRegistry();
  const scheduler = resolveScheduler(schedulerRegistry, requestedId);
  return { schedulerRegistry, scheduler };
}

export function useSchedulerContext() {
  return useRouteContext({ strict: false }) as Partial<
    Awaited<ReturnType<typeof schedulerContext>>
  >;
}

export function useScheduler() {
  const context = useSchedulerContext();
  if (!context.scheduler) {
    throw new Error("Scheduler context is unavailable");
  }
  return {
    scheduler: context.scheduler,
    client: getBuildQueueStateClient(context.scheduler.id),
  };
}
