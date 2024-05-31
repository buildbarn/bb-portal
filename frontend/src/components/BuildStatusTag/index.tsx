import React from 'react';
import {
  CheckCircleFilled,
  CloseCircleFilled,
  InfoCircleFilled,
  LoadingOutlined,
  MinusCircleFilled,
  QuestionCircleFilled,
  StopOutlined,
} from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

export const ALL_STATUS_VALUES = [
  'IN_PROGRESS',
  'SUCCESS',
  'UNSTABLE',
  'FAILURE',
  'NOT_BUILT',
  'ABORTED',
  'UNKNOWN',
] as const;
export type StatusTuple = typeof ALL_STATUS_VALUES;
export type BuildStatusEnum = StatusTuple[number];

interface Props {
  status: BuildStatusEnum;
}

const STATUS_TAGS: { [key in BuildStatusEnum]: React.ReactNode } = {
  IN_PROGRESS: (
    <Tag icon={<LoadingOutlined />} color="blue" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Building</span>
    </Tag>
  ),
  SUCCESS: (
    <Tag icon={<CheckCircleFilled />} color="green" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Succeeded</span>
    </Tag>
  ),
  UNSTABLE: (
    <Tag icon={<InfoCircleFilled />} color="orange" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Unstable</span>
    </Tag>
  ),
  FAILURE: (
    <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Failed</span>
    </Tag>
  ),
  NOT_BUILT: (
    <Tag icon={<MinusCircleFilled />} color="purple" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Not Built</span>
    </Tag>
  ),
  ABORTED: (
    <Tag icon={<StopOutlined />} color="gold" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Aborted</span>
    </Tag>
  ),
  UNKNOWN: (
    <Tag icon={<QuestionCircleFilled />} className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Status Unknown</span>
    </Tag>
  ),
};

const BuildStatusTag: React.FC<Props> = ({ status }) => {
  const resultTag = STATUS_TAGS[status] || STATUS_TAGS.UNKNOWN;
  return <>{resultTag}</>;
};

export default BuildStatusTag;
