function flattenObject(
  obj: unknown,
  path: string,
  out: Record<string, string>,
) {
  if (obj === null || obj === undefined) {
    out[path] = "?";
  } else if (typeof obj === "object") {
    const record = obj as Record<string, unknown>;
    for (const key of Object.keys(record)) {
      const newPath = path ? `${path}/${key}` : key;
      flattenObject(record[key], newPath, out);
    }
  } else {
    out[path] = String(obj);
  }
}

export function flattenUserInfo(obj: unknown): Record<string, string> {
  const out: Record<string, string> = {};
  flattenObject(obj, "", out);
  return out;
}
