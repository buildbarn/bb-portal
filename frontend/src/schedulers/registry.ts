// A scheduler is identified by the instance name prefix it serves, the same
// key bb_storage uses to route Execute requests to it.
export interface Scheduler {
  instanceNamePrefix: string;
}

export interface SchedulerRegistry {
  schedulers: Scheduler[];
}

// "required": several schedulers exist and none was selected. There is no
// default, so the user chooses which prefix to show.
export type SchedulerSelectionProblem =
  | "none-configured"
  | "required"
  | "unknown";

export class SchedulerSelectionError extends Error {
  constructor(
    message: string,
    public registry: SchedulerRegistry,
    public problem: SchedulerSelectionProblem,
  ) {
    super(message);
  }
}

let registryRequest: Promise<SchedulerRegistry> | undefined;

export function loadSchedulerRegistry(): Promise<SchedulerRegistry> {
  registryRequest ??= fetch("/api/v1/schedulers")
    .then(async (response) => {
      if (!response.ok) {
        throw new Error(`Unable to load schedulers: HTTP ${response.status}`);
      }
      return (await response.json()) as SchedulerRegistry;
    })
    .catch((error) => {
      registryRequest = undefined;
      throw error;
    });
  return registryRequest;
}

export function resolveScheduler(
  registry: SchedulerRegistry,
  requested?: unknown,
): Scheduler {
  if (requested === undefined) {
    if (registry.schedulers.length === 1) return registry.schedulers[0];
    throw registry.schedulers.length === 0
      ? new SchedulerSelectionError(
          "No schedulers are configured.",
          registry,
          "none-configured",
        )
      : new SchedulerSelectionError("Select a scheduler", registry, "required");
  }
  if (typeof requested !== "string") {
    throw new SchedulerSelectionError(
      "Invalid scheduler selection: expected an instance name prefix",
      registry,
      "unknown",
    );
  }
  // Exact match only: selections come from the selector or shared links.
  const scheduler = registry.schedulers.find(
    (candidate) => candidate.instanceNamePrefix === requested,
  );
  if (!scheduler) {
    throw new SchedulerSelectionError(
      `Unknown scheduler: ${requested}`,
      registry,
      "unknown",
    );
  }
  return scheduler;
}

export function schedulerLabel({ instanceNamePrefix }: Scheduler): string {
  return instanceNamePrefix === ""
    ? "(all instance names)"
    : instanceNamePrefix;
}

// Switching intentionally discards operation IDs, queue names, filters and cursors.
export function schedulerSwitchTarget(pathname: string, scheduler: string) {
  return {
    to: pathname.startsWith("/operations") ? "/operations" : "/scheduler",
    search: { scheduler },
  } as const;
}

export function schedulerSearch(search: Record<string, unknown>): {
  scheduler?: unknown;
} {
  if (search.scheduler === undefined) return {};
  return {
    scheduler: search.scheduler,
  };
}
