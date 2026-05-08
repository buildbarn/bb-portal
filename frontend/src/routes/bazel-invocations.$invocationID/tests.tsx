import { createFileRoute } from "@tanstack/react-router";
import { TestTab } from "@/components/TestTab";
import { env } from "@/utils/env";
import { requireFeature } from "@/utils/featureGuard";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/bazel-invocations/$invocationID/tests")({
  component: RouteComponent,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  // TODO: Add backend integration test for this when the loader is implemented
  head: (_ctx) => ({
    meta: [
      {
        title: generatePageTitle([
          "Invocation",
          "Tests",
          _ctx.params.invocationID,
        ]),
      },
    ],
  }),
});

function RouteComponent() {
  const { invocationID } = Route.useParams();
  // TODO (isakstenstrom): Fetch data in the data loader
  return <TestTab invocationId={invocationID} />;
}
