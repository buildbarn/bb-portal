export interface Scheduler {
  id: string;
  label: string;
}

export interface SchedulerRegistry {
  schedulers: Scheduler[];
  defaultSchedulerId: string;
}

export class SchedulerSelectionError extends Error {
  constructor(
    message: string,
    public registry: SchedulerRegistry,
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
  requestedId?: unknown,
): Scheduler {
  if (requestedId !== undefined && typeof requestedId !== "string") {
    throw new SchedulerSelectionError(
      "Invalid scheduler selection: expected a scheduler ID",
      registry,
    );
  }
  const id = requestedId ?? registry.defaultSchedulerId;
  const scheduler = registry.schedulers.find(
    (candidate) => candidate.id === id,
  );
  if (!scheduler) {
    throw new SchedulerSelectionError(
      requestedId === undefined && registry.schedulers.length === 0
        ? "No schedulers are configured."
        : `Unknown scheduler: ${id}`,
      registry,
    );
  }
  return scheduler;
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
