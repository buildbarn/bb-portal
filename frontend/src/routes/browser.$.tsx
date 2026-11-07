import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { BrowserPage } from "@/components/pages/Browser";
import { BrowserWelcomePage } from "@/components/pages/Browser/BrowserWelcome";
import { NotFoundError } from "@/main";
import { BrowserPageSchema } from "@/types/BrowserPageParams";
import { generatePageTitle } from "@/utils/generatePageTitle";

const FileSystemAccessProfileSchema = z.object({
  digest: z
    .object({
      hash: z.string(),
      sizeBytes: z.string(),
    })
    .or(z.undefined()),
  pathHashesBaseHash: z.string(),
});

const ComparedActionSchema = BrowserPageSchema.omit({
  browserPageType: true,
});

const BrowserSearchSchema = z.object({
  fileSystemAccessProfile: FileSystemAccessProfileSchema.optional(),
  comparedAction: ComparedActionSchema.optional(),
  openDirs: z.string().optional(),
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
    const { browserPageUrlShape, validateBrowserPageSlug } = await import(
      "@/types/BrowserPageParams"
    );
    const result = validateBrowserPageSlug((params._splat || "").split("/"));
    if (result.success || result.incomplete) {
      return result;
    }
    throw new NotFoundError(
      "page",
      `${result.error}. Expected ${browserPageUrlShape}.`,
    );
  },
  head: (_ctx) => {
    const result = _ctx.loaderData;
    if (result === undefined) {
      return {
        meta: [{ title: generatePageTitle(["Browser", "Page Not Found"]) }],
      };
    }
    if (!result.success) {
      return { meta: [{ title: generatePageTitle(["Browser"]) }] };
    }
    const browserPageParams = result.data;

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
  const result = Route.useLoaderData();
  if (!result.success) {
    return (
      <BrowserWelcomePage
        note={`${result.error}. Complete the URL: ${result.suggestion}`}
      />
    );
  }
  return <BrowserPage params={result.data} search={search} />;
}
