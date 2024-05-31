import dayjs from 'dayjs';

const preciseTo = (from: dayjs.Dayjs, to: dayjs.Dayjs) => {
  const duration = dayjs.duration(to.diff(from));
  return `${Math.floor(duration.asHours())}:${duration.format('mm:ss')}`;
};

export const humanFriendlyAgo = (timestamp: string) => {
  const duration = dayjs.duration(dayjs(timestamp).diff(dayjs()));
  return duration.humanize(true);
};

export default preciseTo;
