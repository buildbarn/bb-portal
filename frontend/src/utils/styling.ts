import type { CSSProperties } from "react";

export function mergeStyles(
  ...styles: (CSSProperties | boolean | undefined)[]
) {
  return styles
    .filter((a) => typeof a === "object")
    .reduce<CSSProperties>((a, b) => Object.assign(a, b), {});
}
