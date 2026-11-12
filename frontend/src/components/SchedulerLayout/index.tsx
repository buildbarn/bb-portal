import { Outlet, useNavigate, useRouterState } from "@tanstack/react-router";
import { Alert, Select, Space } from "antd";
import { DisabledPage } from "@/components/pages/DisabledPage";
import { useSchedulerContext } from "@/schedulers/context";
import {
  type SchedulerRegistry,
  SchedulerSelectionError,
  schedulerSwitchTarget,
} from "@/schedulers/registry";
import { FeatureDisabledError } from "@/utils/featureGuard";

function Selector({
  registry,
  selected,
}: {
  registry: SchedulerRegistry;
  selected?: string;
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
      options={registry.schedulers.map(({ id, label }) => ({
        value: id,
        label,
      }))}
      onChange={(id) => {
        void navigate(schedulerSwitchTarget(pathname, id));
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
        <Selector registry={schedulerRegistry} selected={scheduler.id} />
      )}
      <Outlet key={scheduler.id} />
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
  return (
    <Space orientation="vertical" style={{ width: "100%" }}>
      <Alert
        type="error"
        showIcon
        title="Unable to open scheduler"
        description={
          error instanceof SchedulerSelectionError
            ? registry?.schedulers.length === 0
              ? "No schedulers are configured."
              : "The selected scheduler is not configured. Choose a scheduler below."
            : "Unable to load this scheduler. Please try again or choose another scheduler."
        }
      />
      {registry && registry.schedulers.length > 0 && (
        <Selector registry={registry} />
      )}
    </Space>
  );
}
