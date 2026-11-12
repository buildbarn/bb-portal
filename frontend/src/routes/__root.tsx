import {
  createRootRoute,
  HeadContent,
  retainSearchParams,
} from "@tanstack/react-router";
import { RootLayout } from "@/components/pages/RootLayout";
import { schedulerSearch } from "@/schedulers/registry";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createRootRoute({
  validateSearch: schedulerSearch,
  search: { middlewares: [retainSearchParams(["scheduler"])] },
  component: () => (
    <>
      <HeadContent />
      <RootLayout />
    </>
  ),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle([]) }] }),
});
