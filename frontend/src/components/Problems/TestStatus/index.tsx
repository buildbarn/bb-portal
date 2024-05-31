import React from 'react';
import {
  CheckCircleFilled,
  ClockCircleFilled,
  CloseCircleFilled,
  QuestionCircleFilled,
  StopFilled,
  WarningFilled,
} from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

interface TestStatusProps {
  label: string;
  icon: React.ReactNode; // Not sure if there's a better type we can use?
  color: string;
}

export const TestStatusMap: Map<string, TestStatusProps> = new Map([
  [
    'NO_STATUS',
    {
      label: 'No Status',
      icon: <QuestionCircleFilled />,
      color: 'default',
    },
  ],
  [
    'PASSED',
    {
      label: 'Passed',
      icon: <CheckCircleFilled />,
      color: 'green',
    },
  ],
  [
    'FLAKY',
    {
      label: 'Flaky',
      icon: <WarningFilled />,
      color: 'yellow',
    },
  ],
  [
    'TIMEOUT',
    {
      label: 'Timeout',
      icon: <ClockCircleFilled />,
      color: 'blue',
    },
  ],
  [
    'FAILED',
    {
      label: 'Failed',
      icon: <CloseCircleFilled />,
      color: 'red',
    },
  ],
  [
    'INCOMPLETE',
    {
      label: 'Incomplete',
      icon: <StopFilled />,
      color: 'default',
    },
  ],
  [
    'REMOTE_FAILURE',
    {
      label: 'Remote Failure',
      icon: <CloseCircleFilled />,
      color: 'red',
    },
  ],
  [
    'FAILED_TO_BUILD',
    {
      label: 'Failed to Build',
      icon: <CloseCircleFilled />,
      color: 'red',
    },
  ],
  [
    'TOOL_HALTED_BEFORE_TESTING',
    {
      label: 'Tool Halted Before Testing',
      icon: <StopFilled />,
      color: 'default',
    },
  ],
]);

interface Props {
  status: string;
}

export const TestStatus: React.FC<Props> = ({ status }) => {
  const mappedStatus = TestStatusMap.get(status);
  if (!mappedStatus) {
    return <Tag className={themeStyles.tag}>{status}</Tag>;
  }
  return (
    <Tag className={themeStyles.tag} icon={mappedStatus.icon} color={mappedStatus.color}>
      {mappedStatus.label}
    </Tag>
  );
};
