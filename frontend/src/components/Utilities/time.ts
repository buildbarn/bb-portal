import type { Duration } from "@/lib/grpc-client/google/protobuf/duration";
import dayjs from "dayjs";

export const durationToSeconds = (duration: Duration) => {
  return Number.parseInt(duration.seconds) + duration.nanos / 1e9;
};

const preciseTo = (from: dayjs.Dayjs, to: dayjs.Dayjs) => {
  const duration = dayjs.duration(to.diff(from));
  return `${Math.floor(duration.asHours())}:${duration.format("mm:ss")}`;
};

export const humanFriendlyAgo = (timestamp: string) => {
  const duration = dayjs.duration(dayjs(timestamp).diff(dayjs()));
  return duration.humanize(true);
};

export function millisecondsToTime(milliseconds: number): string {
  const totalSeconds = Math.floor(milliseconds / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  const remainingMilliseconds = Math.floor(milliseconds % 1000);

  return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}:${pad(remainingMilliseconds, 3)}`;
}

function pad(num: number, size = 2): string {
  return num.toString().padStart(size, "0");
}
export default preciseTo;
