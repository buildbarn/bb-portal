import { ClockCircleFilled } from "@ant-design/icons";
import { Divider, Popover, Tag } from "antd";
import type React from "react";
import { useEffect, useState } from "react";
import dayjs from "@/lib/dayjs";
import themeStyles from "@/theme/theme.module.css";
import { readableDurationFromDates } from "@/utils/time";
import styles from "./index.module.css";

interface Props {
  from?: string | null;
  to?: string | null;
  includeIcon?: boolean;
  includePopover?: boolean;
}

const PortalDuration: React.FC<Props> = ({
  from,
  to,
  includeIcon,
  includePopover,
}) => {
  const [now, setNow] = useState(new Date());

  useEffect(() => {
    if (!to) {
      const intervalID = setInterval(() => {
        setNow(new Date());
      }, 500);
      return () => clearInterval(intervalID);
    }
  }, [to]);

  if (!from) return "Unknown";

  const actualFrom = new Date(from);
  const actualTo = !to ? now : new Date(to);
  const content = (
    <Tag
      icon={includeIcon && <ClockCircleFilled />}
      bordered={false}
      className={themeStyles.tagClickable}
    >
      <div className={styles.duration}>
        {readableDurationFromDates(actualFrom, actualTo, { smallestUnit: "s" })}
      </div>
    </Tag>
  );
  if (!includePopover) return content;
  return (
    <Popover
      trigger="click"
      content={
        <div className={styles.popover}>
          {dayjs(actualFrom).format("dddd, MMMM Do, YYYY, [at] h:mm:ss A z")}
          <Divider className={styles.divider}>&darr;</Divider>
          {to
            ? dayjs(actualTo).format("dddd, MMMM Do, YYYY, [at] h:mm:ss A z")
            : "Present"}
        </div>
      }
    >
      {content}
    </Popover>
  );
};

export default PortalDuration;
