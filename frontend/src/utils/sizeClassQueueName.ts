import { z } from "zod";

const PlatformPropertySchema = z.object({
  name: z.string(),
  value: z.string(),
});

export const SizeClassQueueNameSchema = z.object({
  platformQueueName: z.object({
    instanceNamePrefix: z.string(),
    platform: z.object({
      properties: z.array(PlatformPropertySchema),
    }),
  }),
  sizeClass: z.number().int().nonnegative().default(0).catch(0),
});

export type SizeClassQueueNameSearch = z.infer<typeof SizeClassQueueNameSchema>;
