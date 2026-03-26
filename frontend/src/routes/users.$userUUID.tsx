import { createFileRoute } from "@tanstack/react-router";
import { UserDetailsPage } from "@/components/pages/UserDetails";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/users/$userUUID")({
  component: RouteComponent,
  head: (_ctx) => ({
    meta: [{ title: generatePageTitle(["User", _ctx.params.userUUID]) }],
  }),
});

function RouteComponent() {
  const { userUUID } = Route.useParams();
  return <UserDetailsPage userUUID={userUUID} />;
}
