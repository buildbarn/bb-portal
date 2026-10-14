import { z } from "zod";
import type { InvocationName } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { Any } from "@/lib/grpc-client/google/protobuf/any";
import { protobufToObjectWithTypeField } from "@/utils/protobufToObject";

export const PlatformProperySchema = z.object({
  name: z.string(),
  value: z.string(),
});

export const SizeClassQueueNameSchema = z.object({
  platformQueueName: z.object({
    instanceNamePrefix: z.string(),
    platform: z.object({
      properties: z.array(PlatformProperySchema),
    }),
  }),
  sizeClass: z.number().int().nonnegative().default(0).catch(0),
});

export const InvocationIdSchema = z.object({
  typeUrl: z.string(),
  value: z.string(),
});

export const InvocationNameSchema = z.object({
  sizeClassQueueName: SizeClassQueueNameSchema,
  ids: z.array(InvocationIdSchema).default([]),
});

export type SizeClassQueueNameSearch = z.infer<typeof SizeClassQueueNameSchema>;
export type InvocationIdSearch = z.infer<typeof InvocationIdSchema>;
export type InvocationNameSearch = z.infer<typeof InvocationNameSchema>;

export const anyToInvocationIdSearch = (id: Any): InvocationIdSearch => ({
  typeUrl: id.typeUrl,
  value: btoa(Array.from(id.value, (b) => String.fromCharCode(b)).join("")),
});

export const invocationIdSearchToAny = (id: InvocationIdSearch): Any => ({
  typeUrl: id.typeUrl,
  value: Uint8Array.from(atob(id.value), (c) => c.charCodeAt(0)),
});

export const invocationNameFromSearch = (
  search: InvocationNameSearch,
): InvocationName => ({
  sizeClassQueueName: search.sizeClassQueueName,
  ids: search.ids.map(invocationIdSearchToAny),
});

export const rootInvocationNameSearch = (
  sizeClassQueueName: SizeClassQueueNameSearch,
): InvocationNameSearch => ({ sizeClassQueueName, ids: [] });

export const childInvocationNameSearch = (
  parent: InvocationNameSearch,
  childId: Any,
): InvocationNameSearch => ({
  sizeClassQueueName: parent.sizeClassQueueName,
  ids: [...parent.ids, anyToInvocationIdSearch(childId)],
});

export const invocationIdToString = (id: Any): string => {
  const object = protobufToObjectWithTypeField(id, false);
  if (
    typeof object !== "object" ||
    object === null ||
    Object.keys(object).length === 0
  ) {
    return JSON.stringify({ "@type": id.typeUrl });
  }
  return JSON.stringify(object);
};
