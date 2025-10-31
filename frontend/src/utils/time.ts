import type { Duration } from "dayjs/plugin/duration";
import dayjs from "@/lib/dayjs";
import type { Duration as ProtobufDuration } from "@/lib/grpc-client/google/protobuf/duration";

export interface ReadableFormatConfig {
  precision?: number;
  smallestUnit?: "d" | "h" | "m" | "s" | "ms" | "us" | "ns";
}

export const readableDuration = (
  duration: Duration,
  formatConfig: ReadableFormatConfig = {},
): string => {
  const precision = formatConfig.precision ?? 3;
  const smallestUnit = formatConfig.smallestUnit;

  if (duration.asDays() >= 1 || smallestUnit === "d") {
    if (smallestUnit === "d") {
      return `${duration.days()}d`;
    }
    return `${duration.days()}d ${duration.hours()}h`;
  }

  if (duration.asHours() >= 1 || smallestUnit === "h") {
    if (smallestUnit === "h") {
      return `${duration.hours()}h`;
    }
    return `${duration.hours()}h ${duration.minutes()}m`;
  }

  if (duration.asMinutes() >= 2 || smallestUnit === "m") {
    if (smallestUnit === "m") {
      return `${duration.minutes()}m`;
    }
    return `${duration.minutes()}m ${duration.seconds()}s`;
  }

  if (duration.asSeconds() >= 1 || smallestUnit === "s") {
    if (smallestUnit === "s") {
      return `${Math.floor(duration.asSeconds())}s`;
    }
    return `${duration.asSeconds().toPrecision(precision)}s`;
  }

  if (duration.asMilliseconds() >= 1 || smallestUnit === "ms") {
    if (smallestUnit === "ms") {
      return `${duration.milliseconds()}ms`;
    }
    return `${duration.asMilliseconds().toPrecision(precision)}ms`;
  }

  const microseconds = duration.asMilliseconds() * 1000;
  if (microseconds >= 1 || smallestUnit === "us") {
    if (smallestUnit === "us") {
      return `${Math.floor(microseconds)}us`;
    }
    return `${microseconds.toPrecision(precision)}us`;
  }

  const nanoseconds = duration.asMilliseconds() * 1e6;
  if (nanoseconds >= 1 || smallestUnit === "ns") {
    if (smallestUnit === "ns") {
      return `${Math.floor(nanoseconds)}ns`;
    }
    return `${nanoseconds.toPrecision(precision)}ns`;
  }
  return "0ns";
};

export const readableDurationFromSeconds = (
  totalSeconds: number,
  formatConfig: ReadableFormatConfig = {},
): string => {
  return readableDuration(
    dayjs.duration(totalSeconds, "seconds"),
    formatConfig,
  );
};

export const readableDurationFromMilliseconds = (
  totalMilliseconds: number,
  formatConfig: ReadableFormatConfig = {},
): string => {
  return readableDuration(
    dayjs.duration(totalMilliseconds, "milliseconds"),
    formatConfig,
  );
};

export const readableDurationFromProtobufDuration = (
  duration: ProtobufDuration,
  formatConfig: ReadableFormatConfig = {},
): string => {
  return readableDurationFromSeconds(
    protobufDurationToSeconds(duration),
    formatConfig,
  );
};

export const readableDurationFromDates = (
  start: Date,
  end: Date,
  formatConfig: ReadableFormatConfig = {},
): string => {
  const deltaMilliseconds = end.getTime() - start.getTime();
  return readableDurationFromMilliseconds(deltaMilliseconds, formatConfig);
};

export const protobufDurationToSeconds = (duration: ProtobufDuration) => {
  return Number.parseInt(duration.seconds) + duration.nanos / 1e9;
};
