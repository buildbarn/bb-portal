import { createFileRoute } from '@tanstack/react-router';
import { BuildDetailsPage } from '@/components/pages/BuildDetails';
import { generatePageTitle } from '@/utils/generatePageTitle';

export const Route = createFileRoute('/builds/$buildUUID')({
  component: RouteComponent,
  head: (_ctx) => ({meta: [{title: generatePageTitle(["Build", _ctx.params.buildUUID])}]})
});

function RouteComponent() {
  const { buildUUID } = Route.useParams();
  return <BuildDetailsPage buildUUID={buildUUID} />;
}
