import { z } from "zod";
import { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import {
  BrowserPageType,
  getBrowserPageTypeFromString,
} from "@/types/BrowserPageType";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";

export const BrowserPageSchema = z.object({
  instanceName: z.string(),
  digestFunction: z.enum(DigestFunction_Value),
  browserPageType: z.enum(BrowserPageType),
  digest: z.object({
    hash: z
      .string()
      .regex(
        /^[a-f0-9]{64}$/,
        "hash must only contain hex digits and be 64 long",
      ),
    sizeBytes: z.string().regex(/^\d+$/, "sizeBytes must only contain digits"),
  }),
});

export type BrowserPageParams = z.infer<typeof BrowserPageSchema>;

export const parseBrowserPageSlug = (
  slug: Array<string>,
): BrowserPageParams | undefined => {
  const blobIndex = slug.indexOf("blobs");
  if (blobIndex === -1 || blobIndex + 3 >= slug.length) {
    return undefined;
  }

  const instanceName = slug.slice(0, blobIndex).join("/");
  const digestFunction = digestFunctionValueFromString(slug[blobIndex + 1]);
  const browserPageType = getBrowserPageTypeFromString(slug[blobIndex + 2]);

  if (digestFunction === undefined || browserPageType === undefined) {
    return undefined;
  }

  const hashAndSize = slug[blobIndex + 3];
  const [hash, sizeBytes] = hashAndSize.split("-");

  if (!hash || !sizeBytes) {
    return undefined;
  }

  try {
    return BrowserPageSchema.parse({
      instanceName,
      digestFunction,
      browserPageType,
      digest: { hash, sizeBytes },
    });
  } catch (_e) {
    return undefined;
  }
};
