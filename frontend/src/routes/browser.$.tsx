import { createFileRoute, notFound } from '@tanstack/react-router';
import { parseBrowserPageSlug } from "@/utils/parseBrowserPageSlug";
import { generatePageTitle } from '@/utils/generatePageTitle';
import { BrowserPage } from '@/components/pages/Browser';
import { z } from 'zod';

const BrowserSearchSchema = z.object({
  fileSystemAccessProfile: z.object({
    digest: z.object({
      hash: z.string(),
      sizeBytes: z.string(),
    }).or(z.undefined()),
    pathHashesBaseHash: z.string(),
  }).optional(),
})

export type BrowserSearchParams = z.infer<typeof BrowserSearchSchema>

export const Route = createFileRoute('/browser/$')({
  component: RouteComponent,
  validateSearch: (search) => BrowserSearchSchema.parse(search),
  loader: ({ params }) => {
    const browserPageParams = parseBrowserPageSlug((params._splat || '').split('/'));
    if (!browserPageParams) {
      throw notFound();
    }
    return browserPageParams
  },
  head: (_ctx) => {
    const browserPageParams = _ctx.loaderData
    if (browserPageParams === undefined) {
      return { meta: [{ title: generatePageTitle(["Browser", "Page Not Found"]) }] };
    }

    const pageName = browserPageParams.browserPageType
      .split("_")
      .map((word) => word[0].toLocaleUpperCase() + word.slice(1))
      .join(" ");

    return {
      meta: [{
        title: generatePageTitle(["Browser", pageName, browserPageParams.digest.hash])
      }]
    };
  },
});

function RouteComponent() {
  const search = Route.useSearch();
  const browserPageParams = Route.useLoaderData();
  return <BrowserPage params={browserPageParams} search={search} />;
}
