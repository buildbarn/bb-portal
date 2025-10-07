export const readableFileSize = (sizeBytes: number, precision = 3): string => {
  if (sizeBytes < 1024) {
    return `${sizeBytes} B`;
  }

  const kb = sizeBytes / 1024;
  if (kb < 1024) {
    return `${kb.toPrecision(precision)}kiB`;
  }

  const mb = kb / 1024;
  if (mb < 1024) {
    return `${mb.toPrecision(precision)}MiB`;
  }

  const gb = mb / 1024;
  if (gb < 1024) {
    return `${gb.toPrecision(precision)}GiB`;
  }

  const tb = gb / 1024;
  return `${tb.toPrecision(precision)}TiB`;
};

export const readableFileSizeFromString = (
  sizeBytes: string,
  precision = 3,
): string => {
  return readableFileSize(Number.parseInt(sizeBytes), precision);
};
