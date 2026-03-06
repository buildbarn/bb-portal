import { createFileRoute } from '@tanstack/react-router';
import { generatePageTitle } from '@/utils/generatePageTitle';
import { BazelInvocationsPage } from '@/components/pages/BazelInvocations';

export const Route = createFileRoute('/bazel-invocations/')({
  component: BazelInvocationsPage,
  head: (_ctx) => ({meta: [{title: generatePageTitle(["Invocations"])}]})
});
