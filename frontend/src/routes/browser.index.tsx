import { createFileRoute } from "@tanstack/react-router";
import { BrowserWelcomePage } from "@/components/pages/Browser/BrowserWelcome";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/browser/")({
  component: RouteComponent,
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Browser"]) }] }),
});

function RouteComponent() {
  return <BrowserWelcomePage />;
}
