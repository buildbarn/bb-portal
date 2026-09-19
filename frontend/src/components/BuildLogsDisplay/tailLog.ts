// Bazel rewrites its progress display by emitting ANSI sequences that erase
// lines it has already sent, and the portal normalizes those away again, so the
// end of a running log can change or shrink and cannot simply be appended to.
// Each poll therefore re-reads this many trailing lines rather than the whole
// log, which also bounds how far back a rewrite may reach.
const REREAD_LINES = 256;

// Lines at the start of a re-read that are compared against the cached log to
// prove that the cached prefix still lines up with what the server has now.
const ANCHOR_LINES = 8;

const countLines = (text: string): number => {
  let lines = 0;
  for (let i = text.indexOf("\n"); i !== -1; i = text.indexOf("\n", i + 1)) {
    lines++;
  }
  return lines;
};

// Offset at which `line` starts, or -1 if the text holds fewer lines than that.
const startOfLine = (text: string, line: number): number => {
  let offset = 0;
  for (let i = 0; i < line; i++) {
    const newline = text.indexOf("\n", offset);
    if (newline === -1) {
      return -1;
    }
    offset = newline + 1;
  }
  return offset;
};

// Reads the log, re-reading only its tail when there is already a copy of it to
// extend. `fetchFrom` must return the log from the given line onwards.
export const fetchTailedLog = async (
  cached: string | undefined,
  fetchFrom: (startLine: number) => Promise<string>,
): Promise<string> => {
  if (cached === undefined) {
    return fetchFrom(0);
  }

  const firstReadLine = countLines(cached) - REREAD_LINES;
  if (firstReadLine <= 0) {
    return fetchFrom(0);
  }

  const keepUpTo = startOfLine(cached, firstReadLine);
  const anchorEnd = startOfLine(cached, firstReadLine + ANCHOR_LINES);
  if (keepUpTo === -1 || anchorEnd === -1) {
    return fetchFrom(0);
  }

  const tail = await fetchFrom(firstReadLine);
  // A mismatch means lines were erased further back than the re-read, so the
  // cached prefix no longer lines up and the whole log has to be read again.
  if (!tail.startsWith(cached.slice(keepUpTo, anchorEnd))) {
    return fetchFrom(0);
  }
  return cached.slice(0, keepUpTo) + tail;
};
