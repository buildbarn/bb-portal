import { createFileRoute } from '@tanstack/react-router';
import { TargetDetailsPage } from '@/components/pages/TargetDetails';
import { generatePageTitle } from '@/utils/generatePageTitle';
import { apolloClient } from '@/components/ApolloWrapper';
import { gql } from '@/graphql/__generated__';
import { requireFeature } from '@/utils/featureGuard';
import { env } from '@/utils/env';
import { TargetNotFoundError } from '@/main';

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


export const Route = createFileRoute('/targets/$targetID/')({
  component: RouteComponent,
  beforeLoad: requireFeature(env.featureFlags?.bes?.pageTargets),
  loader: async ({ params }) => {
    const { data } = await apolloClient.query({
      query: TARGET_METADATA_QUERY,
      variables: { id: params.targetID },
      fetchPolicy: 'cache-first',
    });
    const target = data?.findTargets?.edges?.[0]?.node;
    if (!target) {
      throw new TargetNotFoundError();
    }
    return { target };
  },
  head: (_ctx) => {
    const label = _ctx.loaderData?.target.label
    if (label === undefined) {
      return { meta: [{ title: generatePageTitle(["Target", "Not Found"]) }] };
    }
    return { meta: [{ title: generatePageTitle(["Target", label]) }] }
  },
});

function RouteComponent() {
  const { target } = Route.useLoaderData();
  return <TargetDetailsPage
    aspect={target.aspect}
    instanceName={target.instanceName.name}
    targetKind={target.targetKind}
    label={target.label}
  />
}
