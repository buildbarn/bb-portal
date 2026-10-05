import z from "zod";

const criticalEventSchema = z.object({
  name: z.string(),
  ts: z.number(),
  dur: z.number(),
});

const criticalPathSchema = criticalEventSchema.array();

export type CriticalPath = z.infer<typeof criticalPathSchema>;
export type CriticalEvent = z.infer<typeof criticalEventSchema>;

// Thrown when the profile is too large to parse in the browser. Resolved
// before any fetch occurs, purely from the profile digest.
export class ProfileTooLargeError extends Error {}

const CRITICAL_PATH_MARKER = "critical path component";

// Estimated decompressed size of a gzip-compressed profile.
const GZIP_COMPRESSION_RATIO = 10;
// Maximum decompressed profile size the browser pipeline handles.
const MAX_DECOMPRESSED_PROFILE_SIZE_BYTES = 128 * 1024 * 1024;

export const isProfileTooLarge = (bytes: number, isGzip: boolean): boolean => {
  return (
    bytes >
    MAX_DECOMPRESSED_PROFILE_SIZE_BYTES / (isGzip ? GZIP_COMPRESSION_RATIO : 1)
  );
};

export const extractCriticalPath = (trace: unknown): CriticalPath => {
  const traceEvents =
    (trace as { traceEvents?: { cat?: string }[] } | null)?.traceEvents ?? [];
  return criticalPathSchema.parse(
    traceEvents.filter((event) => event.cat === CRITICAL_PATH_MARKER),
  );
};
