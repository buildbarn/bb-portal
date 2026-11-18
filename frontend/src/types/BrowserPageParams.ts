import { z } from "zod";
import { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import {
  BrowserPageType,
  getBrowserPageTypeFromString,
} from "@/types/BrowserPageType";
import {
  digestFunctionValueFromString,
  digestFunctionValueToString,
} from "@/utils/digestFunctionUtils";

export const BrowserPageSchema = z.object({
  instanceName: z.string(),
  digestFunction: z.enum(DigestFunction_Value),
  browserPageType: z.enum(BrowserPageType),
  digest: z.object({
    hash: z
      .string()
      .regex(
        /^[a-f0-9]{32,128}$/,
        "hash must only contain hex digits and resemble a hash in length",
      ),
    sizeBytes: z.string().regex(/^\d+$/, "sizeBytes must only contain digits"),
  }),
});

export type BrowserPageParams = z.infer<typeof BrowserPageSchema>;

const slugPlaceholders = [
  "blobs",
  "<digest function>",
  "<page type>",
  "<hash>-<sizeBytes>",
];

export const browserPageUrlShape = `/browser/<instance name>/${slugPlaceholders.join("/")}`;

const digestFunctionNames = Object.values(DigestFunction_Value)
  .filter(
    (value): value is DigestFunction_Value =>
      typeof value === "number" &&
      value !== DigestFunction_Value.UNKNOWN &&
      value !== DigestFunction_Value.UNRECOGNIZED,
  )
  .map(digestFunctionValueToString);

// `incomplete` marks a prefix of a valid URL, such as a breadcrumb link.
export type BrowserPageSlugResult =
  | { success: true; data: BrowserPageParams }
  | { success: false; incomplete: false; error: string }
  | { success: false; incomplete: true; error: string; suggestion: string };

export const validateBrowserPageSlug = (
  rawSlug: Array<string>,
): BrowserPageSlugResult => {
  const slug = rawSlug.filter((segment) => segment !== "");
  const blobIndex = slug.indexOf("blobs");
  const incomplete = (error: string): BrowserPageSlugResult => {
    const present = blobIndex === -1 ? 0 : slug.length - blobIndex;
    return {
      success: false,
      incomplete: true,
      error,
      suggestion: `/browser/${[...slug, ...slugPlaceholders.slice(present)].join("/")}`,
    };
  };
  const invalid = (error: string): BrowserPageSlugResult => ({
    success: false,
    incomplete: false,
    error,
  });

  if (blobIndex === -1) {
    return incomplete('Missing "blobs" segment');
  }

  const instanceName = slug.slice(0, blobIndex).join("/");
  const digestFunctionString = slug[blobIndex + 1];
  if (digestFunctionString === undefined) {
    return incomplete("Missing digest function");
  }
  const digestFunction = digestFunctionValueFromString(digestFunctionString);
  if (digestFunction === DigestFunction_Value.UNRECOGNIZED) {
    return invalid(
      `Unknown digest function "${digestFunctionString}", expected one of: ${digestFunctionNames.join(", ")}`,
    );
  }

  const browserPageTypeString = slug[blobIndex + 2];
  if (browserPageTypeString === undefined) {
    return incomplete("Missing page type");
  }
  const browserPageType = getBrowserPageTypeFromString(browserPageTypeString);
  if (browserPageType === undefined) {
    return invalid(
      `Unknown page type "${browserPageTypeString}", expected one of: ${Object.values(BrowserPageType).join(", ")}`,
    );
  }

  const hashAndSize = slug[blobIndex + 3];
  if (hashAndSize === undefined) {
    return incomplete("Missing digest");
  }
  const [hash, sizeBytes] = hashAndSize.split("-");
  if (!hash || !sizeBytes) {
    return invalid(
      `Malformed digest "${hashAndSize}", expected "<hash>-<sizeBytes>"`,
    );
  }

  const result = BrowserPageSchema.safeParse({
    instanceName,
    digestFunction,
    browserPageType,
    digest: { hash, sizeBytes },
  });
  if (!result.success) {
    return invalid(
      `Malformed digest "${hashAndSize}": ${result.error.issues.map((issue) => issue.message).join(", ")}`,
    );
  }
  return { success: true, data: result.data };
};

export const parseBrowserPageSlug = (
  slug: Array<string>,
): BrowserPageParams | undefined => {
  const result = validateBrowserPageSlug(slug);
  return result.success ? result.data : undefined;
};
