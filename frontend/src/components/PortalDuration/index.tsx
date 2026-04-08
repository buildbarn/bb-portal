import { ClockCircleFilled } from "@ant-design/icons";
import { Divider, Popover, Tag } from "antd";
import type React from "react";
import dayjs from "@/lib/dayjs";
import themeStyles from "@/theme/theme.module.css";
import { type ReadableFormatConfig, readableDurationFromDates } from "@/utils/time";
import styles from "./index.module.css";

interface Props {
  from: string | undefined;
  to: string | undefined;
  includeIcon?: boolean;
  includePopover?: boolean;
  formatConfig?: ReadableFormatConfig;
}

const PortalDuration: React.FC<Props> = ({
  from,
  to,
  includeIcon,
  includePopover,
  formatConfig,
}) => {
  const content = (
    <Tag
      icon={includeIcon && <ClockCircleFilled />}
      bordered={false}
      className={themeStyles.tagClickable}
    >
      <div className={styles.duration}>
        {from && to
          ? readableDurationFromDates(
              new Date(from),
              new Date(to),
              formatConfig,
            )
          : "Unknown"}
      </div>
    </Tag>
  );
  if (!includePopover) return content;
  return (
    <Popover
      trigger="click"
      content={
        <div className={styles.popover}>
          {from
            ? dayjs(from).format("dddd, MMMM Do, YYYY, [at] h:mm:ss A z")
            : "Unknown"}
          <Divider className={styles.divider}>&darr;</Divider>
          {to
            ? dayjs(to).format("dddd, MMMM Do, YYYY, [at] h:mm:ss A z")
            : "Unknown"}
        </div>
      }
    >
      {content}
    </Popover>
  );
};

export default PortalDuration;
