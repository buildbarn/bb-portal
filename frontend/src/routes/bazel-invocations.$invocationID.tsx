import { createFileRoute } from "@tanstack/react-router";
import { BazelInvocationDetailsPage } from "@/components/pages/BazelInvocationDetails";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/bazel-invocations/$invocationID")({
  component: RouteComponent,
  head: (_ctx) => ({
    meta: [
      { title: generatePageTitle(["Invocation", _ctx.params.invocationID]) },
    ],
  }),
});

function RouteComponent() {
  const { invocationID } = Route.useParams();
  return <BazelInvocationDetailsPage invocationID={invocationID} />;
}
