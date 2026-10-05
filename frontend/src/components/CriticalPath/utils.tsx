import z from "zod";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";
import { generateFileUrl } from "@/utils/urlGenerator";

const criticalEventSchema = z.object({
  name: z.string(),
  ts: z.number(),
  dur: z.number(),
});

const criticalPathSchema = criticalEventSchema.array();

export type CriticalPath = z.infer<typeof criticalPathSchema>;
export type CriticalEvent = z.infer<typeof criticalEventSchema>;

const fetchProfileBytes = async (
  profile: FileDetailsFragment,
): Promise<Uint8Array> => {
  const res = await fetch(
    generateFileUrl(
      profile.digest.rev2InstanceName,
      digestFunctionValueFromString(profile.digest.digestFunction),
      {
        hash: profile.digest.hash,
        sizeBytes: profile.digest.sizeBytes.toString(),
      },
      profile.filePath.path,
    ),
  );
  if (!res.ok) {
    throw new Error(`Failed to fetch profile: ${res.statusText}`);
  }
  return new Uint8Array(await res.arrayBuffer());
};

const isGzipProfile = (profile: FileDetailsFragment) =>
  profile.filePath.path.endsWith(".gz");

export const getCriticalPath = async (
  profile: FileDetailsFragment,
): Promise<CriticalPath | undefined> => {
  try {
    const bytes = await fetchProfileBytes(profile);
    let jsonData = null;
    if (isGzipProfile(profile)) {
      const compressedStream = new Response(bytes.slice().buffer).body;

      const decompressionStream = new DecompressionStream("gzip");
      const decompressedStream =
        compressedStream?.pipeThrough(decompressionStream);
      const decompressedResponse = new Response(decompressedStream);
      jsonData = await decompressedResponse.json();
    } else {
      // We assume the file is a JSON file if not .gz.
      const text = new TextDecoder().decode(bytes);
      jsonData = JSON.parse(text);
    }

    const parsedData = criticalPathSchema.parse(
      jsonData.traceEvents.filter(
        (event: { cat?: string }) =>
          event.cat && event.cat === "critical path component",
      ),
    );
    return parsedData;
  } catch (error) {
    console.error("Failed to fetch critical path data", error);
    return undefined;
  }
};
