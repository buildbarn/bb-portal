import { generatePageTitle } from '@/utils/generatePageTitle';
import { OperationDetailsPage } from '@/components/pages/OperationDetails';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/operations/$operationID')({
  component: RouteComponent,
  head: (_ctx) => ({meta: [{title: generatePageTitle(["Operation", _ctx.params.operationID])}]})
})

function RouteComponent() {
  const { operationID } = Route.useParams();
  return <OperationDetailsPage operationID={operationID} />;
}
