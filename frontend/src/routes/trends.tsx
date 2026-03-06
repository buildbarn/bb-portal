import { TrendsPage } from "@/components/pages/Trends";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute('/trends')({
  component: TrendsPage,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTrends),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Trends"]) }] }),
})
