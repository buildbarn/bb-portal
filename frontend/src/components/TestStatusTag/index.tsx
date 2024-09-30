import React from 'react';
import {
    CheckCircleFilled,
    CloseCircleFilled,
    InfoCircleFilled,
    MinusCircleFilled,
    QuestionCircleFilled,
    StopOutlined,
} from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

export const ALL_STATUS_VALUES = [
    'NO_STATUS',
    'PASSED',
    'FLAKY',
    'TIMEOUT',
    'FAILED',
    'INCOMPLETE',
    'REMOTE_FAILURE',
    'FAILED_TO_BUILD',
    'TOOL_HALTED_BEFORE_TESTING'
] as const;

export type StatusTuple = typeof ALL_STATUS_VALUES;
export type TestStatusEnum = StatusTuple[number];

interface Props {
    status: TestStatusEnum;
}

const STATUS_TAGS: { [key in TestStatusEnum]: React.ReactNode } = {
    NO_STATUS: (
        <Tag icon={<QuestionCircleFilled />} className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>No Status</span>
        </Tag>
    ),
    PASSED: (
        <Tag icon={<CheckCircleFilled />} color="green" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Passed</span>
        </Tag>
    ),
    FLAKY: (
        <Tag icon={<InfoCircleFilled />} color="orange" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Flaky</span>
        </Tag>
    ),
    FAILED: (
        <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Failed</span>
        </Tag>
    ),
    TIMEOUT: (
        <Tag icon={<MinusCircleFilled />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Timeout</span>
        </Tag>
    ),
    INCOMPLETE: (
        <Tag icon={<StopOutlined />} color="blue" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Incomplete</span>
        </Tag>
    ),
    REMOTE_FAILURE: (
        <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Remote Failure</span>
        </Tag>
    ),
    FAILED_TO_BUILD: (
        <Tag icon={<QuestionCircleFilled />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Failed to Build</span>
        </Tag>
    ),
    TOOL_HALTED_BEFORE_TESTING: (
        <Tag icon={<QuestionCircleFilled />} color="blue" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Status Unknown</span>
        </Tag>
    ),
};

const TestStatusTag: React.FC<Props> = ({ status }) => {
    const resultTag = STATUS_TAGS[status] || STATUS_TAGS.NO_STATUS;
    return <>{resultTag}</>;
};

export default TestStatusTag;
