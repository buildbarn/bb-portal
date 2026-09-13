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
