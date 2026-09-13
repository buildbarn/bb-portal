import z from "zod";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { casByteStreamClient } from "@/grpc/casByteStreamClient";
import { digestFunction_ValueFromJSON } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { fetchCasObject } from "@/utils/fetchCasObject";

const criticalEventSchema = z.object({
  name: z.string(),
  ts: z.number(),
  dur: z.number(),
});

const criticalPathSchema = criticalEventSchema.array();

export type CriticalPath = z.infer<typeof criticalPathSchema>;
export type CriticalEvent = z.infer<typeof criticalEventSchema>;

export const getCriticalPath = async (
  profile: FileDetailsFragment,
): Promise<CriticalPath | undefined> => {
  try {
    const res = await fetchCasObject(
      casByteStreamClient,
      profile.digest.rev2InstanceName,
      digestFunction_ValueFromJSON(profile.digest.digestFunction.toUpperCase()),
      {
        hash: profile.digest.hash,
        sizeBytes: profile.digest.sizeBytes.toString(),
      },
    );
    let jsonData = null;
    if (profile.filePath.path.endsWith(".gz")) {
      const compressedStream = new Response(new Uint8Array(res)).body;

      const decompressionStream = new DecompressionStream("gzip");
      const decompressedStream =
        compressedStream?.pipeThrough(decompressionStream);
      const decompressedResponse = new Response(decompressedStream);
      jsonData = await decompressedResponse.json();
    } else {
      // We assume the file is a JSON file if not .gz.
      const text = new TextDecoder().decode(res);
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
