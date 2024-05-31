import React, { useEffect, useState } from 'react';
import { Divider, Popover, Tag } from 'antd';
import { ClockCircleFilled } from '@ant-design/icons';
import styles from './index.module.css';
import themeStyles from '@/theme/theme.module.css';
import dayjs from '@/lib/dayjs';
import preciseTo from '@/components/Utilities/time';

interface Props {
  from?: string | null;
  to?: string | null;
  includePlus?: boolean;
  includeIcon?: boolean;
  includePopover?: boolean;
}

const PortalDuration: React.FC<Props> = ({ from, to, includePlus, includeIcon, includePopover }) => {
  const [now, setNow] = useState(dayjs().toString());
  useEffect(() => {
    const intervalID = setInterval(() => {
      setNow(dayjs().toString());
    }, 500);
    return () => clearInterval(intervalID);
  }, []);
  if (!from) return 'Unknown';
  const actualTo = !to ? now : to;
  const content = (
    <Tag icon={includeIcon && <ClockCircleFilled />} bordered={false} className={themeStyles.tagClickable}>
      <div className={styles.duration}>
        {includePlus && '+'}
        {preciseTo(dayjs(from), dayjs(actualTo))}
      </div>
    </Tag>
  );
  if (!includePopover) return content;
  return (
    <Popover
      trigger="click"
      content={
        <div className={styles.popover}>
          {dayjs(from).format('dddd, MMMM Do, YYYY, [at] h:mm:ss A z')}
          <Divider className={styles.divider}>&darr;</Divider>
          {to ? dayjs(actualTo).format('dddd, MMMM Do, YYYY, [at] h:mm:ss A z') : 'Present'}
        </div>
      }
    >
      {content}
    </Popover>
  );
};

export default PortalDuration;
