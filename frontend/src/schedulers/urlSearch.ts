import {
  defaultParseSearch,
  defaultStringifySearch,
} from "@tanstack/react-router";

// IDs are opaque strings. TanStack's JSON parser otherwise turns valid IDs such
// as "true", "1e3", or long digit sequences into booleans or rounded numbers.
export function parsePortalSearch(search: string): Record<string, unknown> {
  const parsed: Record<string, unknown> = { ...defaultParseSearch(search) };
  const ids = new URLSearchParams(search).getAll("scheduler");
  if (ids.length === 1) {
    let id = ids[0];
    // Preserve links emitted by the previous default JSON string serializer.
    if (id.startsWith('"')) {
      try {
        const decoded: unknown = JSON.parse(id);
        if (typeof decoded === "string") id = decoded;
      } catch {
        // An invalid explicit selection is reported by normal ID resolution.
      }
    }
    parsed.scheduler = id;
  } else if (ids.length > 1) {
    parsed.scheduler = ids;
  }
  return parsed;
}

export function stringifyPortalSearch(search: Record<string, unknown>): string {
  const { scheduler, ...other } = search;
  const serialized = defaultStringifySearch(other);
  if (scheduler === undefined) return serialized;
  if (typeof scheduler !== "string") return defaultStringifySearch(search);
  return `${serialized}${serialized ? "&" : "?"}scheduler=${encodeURIComponent(scheduler)}`;
}
