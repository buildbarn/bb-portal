import { createFileRoute, notFound } from "@tanstack/react-router";
import { z } from "zod";
import { BrowserPage } from "@/components/pages/Browser";
import { generatePageTitle } from "@/utils/generatePageTitle";

const BrowserSearchSchema = z.object({
  fileSystemAccessProfile: z
    .object({
      digest: z
        .object({
          hash: z.string(),
          sizeBytes: z.string(),
        })
        .or(z.undefined()),
      pathHashesBaseHash: z.string(),
    })
    .optional(),
});

export type BrowserSearchParams = z.infer<typeof BrowserSearchSchema>;

export const Route = createFileRoute("/browser/$")({
  component: RouteComponent,
  validateSearch: (search) => BrowserSearchSchema.parse(search),
  loader: async ({ params }) => {
    // Asynchronous import of parseBrowserPageSlug as it depends on the
    // REv2 grpc client. This prevents the client from being loaded for
    // every route as we only need it when actually loading the browser
    // page.
    const { parseBrowserPageSlug } = await import(
      "@/utils/parseBrowserPageSlug"
    );
    const browserPageParams = parseBrowserPageSlug(
      (params._splat || "").split("/"),
    );
    if (!browserPageParams) {
      throw notFound();
    }
    return browserPageParams;
  },
  head: (_ctx) => {
    const browserPageParams = _ctx.loaderData;
    if (browserPageParams === undefined) {
      return {
        meta: [{ title: generatePageTitle(["Browser", "Page Not Found"]) }],
      };
    }

    const pageName = browserPageParams.browserPageType
      .split("_")
      .map((word) => word[0].toLocaleUpperCase() + word.slice(1))
      .join(" ");

    return {
      meta: [
        {
          title: generatePageTitle([
            "Browser",
            pageName,
            browserPageParams.digest.hash,
          ]),
        },
      ],
    };
  },
});

function RouteComponent() {
  const search = Route.useSearch();
  const browserPageParams = Route.useLoaderData();
  return <BrowserPage params={browserPageParams} search={search} />;
}
