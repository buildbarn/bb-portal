import React from 'react';
import {
    CloseCircleFilled,
    ClockCircleOutlined,
    MinusCircleFilled,
    ExclamationCircleOutlined,
    MinusOutlined,
    StopOutlined,
    WarningOutlined,
    AlertOutlined,
    CloseOutlined,

} from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

export const ALL_STATUS_VALUES = [
    'SKIPPED',
    'USER_INTERRUPTED',
    'TIME_OUT',
    'REMOTE_ENVIRONMENT_FAILURE',
    'INTERNAL',
    "LOADING_FAILURE",
    "ANALYSIS_FAILURE",
    "NO_ANALYZE",
    "NO_BUILD",
    "INCOMPLETE",
    "OUT_OF_MEMORY",
    "NONE",
] as const;

export type AbortReasonsTuple = typeof ALL_STATUS_VALUES;
export type AbortReasonsEnum = AbortReasonsTuple[number];

interface Props {
    reason: AbortReasonsEnum;
}

const STATUS_TAGS: { [key in AbortReasonsEnum]: React.ReactNode } = {
    SKIPPED: (
        <Tag icon={<MinusCircleFilled />} className={themeStyles.tag} color="purple" >
            <span className={themeStyles.tagContent}>Skipped</span>
        </Tag>
    ),
    USER_INTERRUPTED: (
        <Tag icon={<StopOutlined />} color="blue" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>User Interrupted</span>
        </Tag>
    ),
    TIME_OUT: (
        <Tag icon={<ClockCircleOutlined />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Time Out</span>
        </Tag>
    ),
    REMOTE_ENVIRONMENT_FAILURE: (
        <Tag icon={<WarningOutlined />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Remote Environment Failure</span>
        </Tag>
    ),
    INTERNAL: (
        <Tag icon={<ExclamationCircleOutlined />} color="cyan" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Internal</span>
        </Tag>
    ),
    LOADING_FAILURE: (
        <Tag icon={<CloseOutlined />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Loading Failure</span>
        </Tag>
    ),
    ANALYSIS_FAILURE: (
        <Tag icon={<CloseOutlined />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Analysis Failure</span>
        </Tag>
    ),
    NO_ANALYZE: (
        <Tag icon={<MinusCircleFilled />} color="geekblue" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>No Analyze</span>
        </Tag>
    ),
    NO_BUILD: (
        <Tag icon={<MinusCircleFilled />} color="geekblue-inverse" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>No Build</span>
        </Tag>
    ),
    INCOMPLETE: (
        <Tag icon={<CloseCircleFilled />} color="gold" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Incomplete</span>
        </Tag>
    ),
    OUT_OF_MEMORY: (
        <Tag icon={<AlertOutlined />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Out of Memory</span>
        </Tag>
    ),
    NONE: (
        <Tag icon={<MinusOutlined />} color="default" className={themeStyles.tag} />

    ),
};

const TargetAbortReasonTag: React.FC<Props> = ({ reason }) => {
    const resultTag = STATUS_TAGS[reason] || STATUS_TAGS.NONE;
    return <>{resultTag}</>;
};

export default TargetAbortReasonTag;
