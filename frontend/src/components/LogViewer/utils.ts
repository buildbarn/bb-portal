// Escape regex special characters to treat query as plain text.
// Examples:
// "a+b"   -> "a\+b"       (prevents + from meaning one or more)
// "file.*" -> "file\.\*"  (prevents .* from matching any characters)
// "a|b"    -> "a\|b"      (prevents | from acting as an OR operator)
// export const escapedQuery = (query: string) => {
//   return query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
// };
export const escapeRegex = (query: string) => {
  const escapedQuery = query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  return new RegExp(escapedQuery, "gi");
};

// Shared so that repeated searches without a query keep the same identity.
const NO_MATCHES: number[] = [];

// Indices of the lines matching `query`, one entry per occurrence, up to
// `limit` occurrences.
export const findMatchIndices = (
  items: string[],
  query: string,
  limit: number,
): number[] => {
  // An empty query escapes to an empty pattern, which matches at every
  // position of every line rather than nowhere.
  if (query === "") {
    return NO_MATCHES;
  }

  const result: number[] = [];
  const escapedQuery = escapeRegex(query);
  for (let i = 0; i < items.length; i++) {
    const cleanItem = items[i].replace(ansiRegex(), "");

    for (const _ of cleanItem.matchAll(escapedQuery)) {
      result.push(i);

      if (result.length >= limit) {
        return result;
      }
    }
  }

  return result;
};

// Regex for ANSI escape codes (colors, formatting, etc)
// Example: "\x1B[31mError\x1B[0m" ignores the color codes
export const ansiRegex = () => {
  const ansiPattern = "\\x1B\\[[0-9;]*m";
  return new RegExp(ansiPattern, "g");
};

// ANSI escape codes for highlighting
export const ANSI_HIGHLIGHT_START = "\x1B[30;103m";

// ANSI escape code for resetting formatting
export const ANSI_RESET = "\x1B[0m";
