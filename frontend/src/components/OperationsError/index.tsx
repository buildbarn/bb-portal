import { SearchParamError } from "@tanstack/react-router";
import PortalAlert from "@/components/PortalAlert";
import { ErrorPage } from "@/components/pages/ErrorPage";
import { useSchedulerContext } from "@/schedulers/context";

export default function OperationsError({ error }: { error: Error }) {
  const { schedulerRegistry } = useSchedulerContext();
  if (error instanceof SearchParamError) return <ErrorPage error={error} />;
  // The parent SchedulerLayout already renders the selector for this route.
  return (
    <PortalAlert
      showIcon
      type="error"
      title="Unable to load operations"
      description={
        (schedulerRegistry?.schedulers.length ?? 0) > 1
          ? "Please try again or choose another scheduler using the selector above."
          : "Please try again."
      }
    />
  );
}
