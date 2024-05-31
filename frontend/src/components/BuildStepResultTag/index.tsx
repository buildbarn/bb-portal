import React from 'react';
import {
  CheckCircleFilled,
  CloseCircleFilled,
  ExclamationCircleFilled,
  InfoCircleFilled,
  QuestionCircleFilled,
  StopFilled,
} from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';
export enum BuildStepResultEnum {
  'SUCCESS'= 'SUCCESS',
  'UNSTABLE' = 'UNSTABLE',
  'PARSING_FAILURE' = 'PARSING_FAILURE',
  'BUILD_FAILURE' = 'BUILD_FAILURE',
  'TESTS_FAILED' = 'TESTS_FAILED',
  'NOT_BUILT' = 'NOT_BUILT',
  'ABORTED' = 'ABORTED',
  'INTERRUPTED' = 'INTERRUPTED',
  'UNKNOWN' =  'UNKNOWN',
}

interface Props {
  result: BuildStepResultEnum;
}

const RESULT_TAGS: { [key in BuildStepResultEnum]: React.ReactNode } = {
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
  PARSING_FAILURE: (
    <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Parsing Failed</span>
    </Tag>
  ),
  BUILD_FAILURE: (
    <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Build Failed</span>
    </Tag>
  ),
  TESTS_FAILED: (
    <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Tests Failed</span>
    </Tag>
  ),
  NOT_BUILT: (
    <Tag icon={<StopFilled />} color="purple" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Not Built</span>
    </Tag>
  ),
  ABORTED: (
    <Tag icon={<ExclamationCircleFilled />} color="cyan" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Aborted</span>
    </Tag>
  ),
  INTERRUPTED: (
    <Tag icon={<ExclamationCircleFilled />} color="cyan" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Interrupted</span>
    </Tag>
  ),
  UNKNOWN: (
    <Tag icon={<QuestionCircleFilled />} color="default" className={themeStyles.tag}>
      <span className={themeStyles.tagContent}>Status Unknown</span>
    </Tag>
  ),
};

const BuildStepResultTag: React.FC<Props> = ({ result }) => {
  const resultTag = RESULT_TAGS[result] || RESULT_TAGS.UNKNOWN;
  return <>{resultTag}</>;
};

export default BuildStepResultTag;
