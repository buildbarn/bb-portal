import { durationToSeconds } from "@/components/Utilities/time";
import type { Duration } from "@/lib/grpc-client/google/protobuf/duration";

export const formatDurationFromSeconds = (
  totalSeconds: number,
  precision = 8,
  millisecondsDecimals?: number,
): string => {
  if (totalSeconds < 1) {
    const milliseconds = totalSeconds * 1e3;
    if (millisecondsDecimals !== undefined) {
      return `${milliseconds.toFixed(millisecondsDecimals)}ms`;
    }
    return `${milliseconds.toPrecision(precision)}ms`;
  }

  if (totalSeconds < 60) {
    return `${totalSeconds.toPrecision(precision)}s`;
  }

  const totalMinutes = totalSeconds / 60;
  if (totalMinutes < 60) {
    return `${totalMinutes.toPrecision(precision)}m`;
  }

  const totalHours = totalMinutes / 60;
  if (totalHours < 24) {
    return `${totalHours.toPrecision(precision)}h`;
  }

  const totalDays = totalHours / 24;
  return `${totalDays.toPrecision(precision)}d`;
};

export const formatDuration = (duration: Duration, precision = 8): string => {
  const totalSeconds = durationToSeconds(duration);
  return formatDurationFromSeconds(totalSeconds, precision);
};

export const formatDurationFromDates = (
  start: Date,
  end: Date,
  precision = 8,
  millisecondsDecimals?: number,
): string => {
  const deltaMilliseconds = end.getTime() - start.getTime();
  return formatDurationFromSeconds(
    deltaMilliseconds / 1000,
    precision,
    millisecondsDecimals,
  );
};

export const formatFileSize = (sizeBytes: number): string => {
  if (sizeBytes < 1024) {
    return `${sizeBytes} B`;
  }

  const kb = sizeBytes / 1024;
  if (kb < 1024) {
    return `${kb.toPrecision(3)} kiB`;
  }

  const mb = kb / 1024;
  if (mb < 1024) {
    return `${mb.toPrecision(3)} MiB`;
  }

  const gb = mb / 1024;
  if (gb < 1024) {
    return `${gb.toPrecision(3)} GiB`;
  }

  const tb = gb / 1024;
  return `${tb.toPrecision(3)} TiB`;
};

export const formatFileSizeFromString = (sizeBytes: string): string => {
  return formatFileSize(Number.parseInt(sizeBytes));
};
