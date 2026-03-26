import { createFileRoute, Outlet } from "@tanstack/react-router";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";

export const Route = createFileRoute("/bazel-invocations")({
  component: Outlet,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageInvocations),
});
