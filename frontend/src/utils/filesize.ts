export const readableFileSize = (
  sizeBytes: number | bigint,
  precision = 3,
): string => {
  const numericSizeBytes = Number(sizeBytes);
  if (numericSizeBytes < 1024) {
    return `${numericSizeBytes} B`;
  }

  const kib = numericSizeBytes / 1024;
  if (kib < 1024) {
    return `${Number(kib.toPrecision(precision))} kiB`;
  }

  const mib = kib / 1024;
  if (mib < 1024) {
    return `${Number(mib.toPrecision(precision))} MiB`;
  }

  const gib = mib / 1024;
  if (gib < 1024) {
    return `${Number(gib.toPrecision(precision))} GiB`;
  }

  const tib = gib / 1024;
  return `${Number(tib.toPrecision(precision))} TiB`;
};

export const readableFileSizeFromString = (
  sizeBytes: string,
  precision = 3,
): string => {
  return readableFileSize(Number.parseInt(sizeBytes, 10), precision);
};
