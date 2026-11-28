import { Outlet, useNavigate, useRouterState } from "@tanstack/react-router";
import { Alert, Select, Space } from "antd";
import { DisabledPage } from "@/components/pages/DisabledPage";
import { useSchedulerContext } from "@/schedulers/context";
import {
  type SchedulerRegistry,
  SchedulerSelectionError,
  schedulerLabel,
  schedulerSwitchTarget,
} from "@/schedulers/registry";
import { FeatureDisabledError } from "@/utils/featureGuard";

function Selector({
  registry,
  selected,
  preserveCurrentLocation = false,
}: {
  registry: SchedulerRegistry;
  selected?: string;
  preserveCurrentLocation?: boolean;
}) {
  const navigate = useNavigate();
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  return (
    <Select
      aria-label="Scheduler"
      placeholder="Select a scheduler"
      value={selected}
      style={{ minWidth: 200 }}
      options={registry.schedulers.map((scheduler) => ({
        value: scheduler.instanceNamePrefix,
        label: schedulerLabel(scheduler),
      }))}
      onChange={(instanceNamePrefix) => {
        if (preserveCurrentLocation) {
          void navigate({
            to: pathname,
            search: (prev) => ({ ...prev, scheduler: instanceNamePrefix }),
          });
        } else {
          void navigate(schedulerSwitchTarget(pathname, instanceNamePrefix));
        }
      }}
    />
  );
}

export function SchedulerLayout() {
  const { schedulerRegistry, scheduler } = useSchedulerContext();
  if (!schedulerRegistry || !scheduler) return null;
  return (
    <Space orientation="vertical" style={{ width: "100%" }}>
      {schedulerRegistry.schedulers.length > 1 && (
        <Selector
          registry={schedulerRegistry}
          selected={scheduler.instanceNamePrefix}
        />
      )}
      <Outlet key={scheduler.instanceNamePrefix} />
    </Space>
  );
}

export function SchedulerError({ error }: { error: Error }) {
  const context = useSchedulerContext();
  const registry =
    error instanceof SchedulerSelectionError
      ? error.registry
      : context.schedulerRegistry;
  if (error instanceof FeatureDisabledError) return <DisabledPage />;
  const isUnknownSelection =
    error instanceof SchedulerSelectionError && error.problem === "unknown";
  const selector =
    registry && (registry.schedulers.length > 1 || isUnknownSelection) ? (
      <Selector
        registry={registry}
        preserveCurrentLocation={
          error instanceof SchedulerSelectionError &&
          error.problem === "required"
        }
      />
    ) : null;
  if (
    error instanceof SchedulerSelectionError &&
    error.problem === "required"
  ) {
    return (
      <Space orientation="vertical" style={{ width: "100%" }}>
        <Alert
          type="info"
          showIcon
          title="Choose a scheduler"
          description="Several schedulers are configured. Choose the instance name prefix to show."
        />
        {selector}
      </Space>
    );
  }
  return (
    <Space orientation="vertical" style={{ width: "100%" }}>
      <Alert
        type="error"
        showIcon
        title="Unable to open scheduler"
        description={
          error instanceof SchedulerSelectionError
            ? error.problem === "none-configured"
              ? "No schedulers are configured."
              : "The selected scheduler is not configured. Choose a scheduler below."
            : selector
              ? "Unable to load this scheduler. Please try again or choose another scheduler."
              : "Unable to load this scheduler. Please try again."
        }
      />
      {selector}
    </Space>
  );
}
