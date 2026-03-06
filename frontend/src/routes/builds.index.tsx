import { createFileRoute } from '@tanstack/react-router';
import { BuildsPage } from '@/components/pages/Builds';
import { generatePageTitle } from '@/utils/generatePageTitle';

export const Route = createFileRoute('/builds/')({
  component: RouteComponent,
  head: (_ctx) => ({meta: [{title: generatePageTitle(["Builds"])}]})
});

function RouteComponent() {
  return <BuildsPage />;
}
