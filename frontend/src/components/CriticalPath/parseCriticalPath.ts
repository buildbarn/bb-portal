import z from "zod";

const criticalEventSchema = z.object({
  name: z.string(),
  ts: z.number(),
  dur: z.number(),
});

const criticalPathSchema = criticalEventSchema.array();

export type CriticalPath = z.infer<typeof criticalPathSchema>;
export type CriticalEvent = z.infer<typeof criticalEventSchema>;

export const extractCriticalPath = (trace: unknown): CriticalPath => {
  const traceEvents =
    (trace as { traceEvents?: { cat?: string }[] } | null)?.traceEvents ?? [];
  return criticalPathSchema.parse(
    traceEvents.filter((event) => event.cat === "critical path component"),
  );
};
