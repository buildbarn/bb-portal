import { generatePageTitle } from '@/utils/generatePageTitle';
import { TestDetailsPage } from '@/components/pages/TestDetails';
import { requireFeature } from '@/utils/featureGuard';
import { createFileRoute } from '@tanstack/react-router';
import { gql } from '@/graphql/__generated__';
import { apolloClient } from '@/components/ApolloWrapper';
import { env } from '@/utils/env';
import { TestNotFoundError } from '@/main';

export const TARGET_METADATA_QUERY = gql(`
  query getTargetMetadata($id: ID!) {
    findTargets(where: { id: $id }) {
      edges {
        node {
          id
          aspect
          instanceName {
            name
          }
          label
          targetKind
        }
      }
    }
  }
`);


export const Route = createFileRoute('/targets/$targetID/tests')({
  component: RouteComponent,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTests),
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: TARGET_METADATA_QUERY,
      variables: { id: params.targetID },
      fetchPolicy: 'cache-first',
    });
    const target = data?.findTargets?.edges?.[0]?.node;
    if (!target) {
      throw new TestNotFoundError();
    }
    return { target };
  },
  head: (_ctx) => {
    const label = _ctx.loaderData?.target.label
    if (label === undefined) {
      return { meta: [{ title: generatePageTitle(["Test", "Not Found"]) }] };
    }
    return { meta: [{ title: generatePageTitle(["Test", label]) }] }
  },});

function RouteComponent() {
  const { target } = Route.useLoaderData();
  return <TestDetailsPage
    targetID={target.id}
    aspect={target.aspect}
    instanceName={target.instanceName.name}
    kind={target.targetKind}
    label={target.label}
  />
}
