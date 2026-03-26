import { createFileRoute } from "@tanstack/react-router";
import { TargetsPage } from "@/components/pages/Targets";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/targets/")({
  component: TargetsPage,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Targets"]) }] }),
});
