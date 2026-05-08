import { ClockCircleFilled } from "@ant-design/icons";
import type React from "react";
import dayjs from "@/lib/dayjs";
import {
  type ReadableFormatConfig,
  readableDurationFromDates,
} from "@/utils/time";
import CodeText from "../CodeText";

interface Props {
  from: string | undefined;
  to: string | undefined;
  includeIcon?: boolean;
  formatConfig?: ReadableFormatConfig;
}

const PortalDuration: React.FC<Props> = ({
  from,
  to,
  includeIcon,
  formatConfig,
}) => {
  const dayjsFormat = "dddd, MMMM Do, YYYY, [at] HH:mm:ss z";
  return (
    <CodeText
      title={`Started at: ${from ? dayjs(from).format(dayjsFormat) : "Unknown"}\nEnded at: ${to ? dayjs(to).format(dayjsFormat) : "Unknown"}`}
    >
      {includeIcon && <ClockCircleFilled style={{ marginRight: "5px" }} />}
      {from && to
        ? readableDurationFromDates(new Date(from), new Date(to), formatConfig)
        : "Unknown"}
    </CodeText>
  );
};

export default PortalDuration;
