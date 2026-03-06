import { createFileRoute } from '@tanstack/react-router';
import { TestsPage } from '@/components/pages/Tests';
import { generatePageTitle } from '@/utils/generatePageTitle';

export const Route = createFileRoute('/tests/')({
  component: TestsPage,
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Tests"]) }] }),
});
