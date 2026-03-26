import { createRootRoute, HeadContent } from "@tanstack/react-router";
import { RootLayout } from "@/components/pages/RootLayout";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createRootRoute({
  component: () => (
    <>
      <HeadContent />
      <RootLayout />
    </>
  ),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle([]) }] }),
});
