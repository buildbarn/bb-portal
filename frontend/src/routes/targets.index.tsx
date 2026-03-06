import { createFileRoute } from '@tanstack/react-router';
import { TargetsPage } from '@/components/pages/Targets';
import { generatePageTitle } from '@/utils/generatePageTitle';
import { requireFeature } from '@/utils/featureGuard';
import { env } from '@/utils/env';

export const Route = createFileRoute('/targets/')({
  component: TargetsPage,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Targets"]) }] }),
});
