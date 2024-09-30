import React from 'react';
import {

    CloseCircleFilled,
    InfoCircleOutlined,
    MinusCircleFilled,
    QuestionCircleFilled,
    StopOutlined,
    FileUnknownOutlined,
    KeyOutlined
} from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

export const ALL_STATUS_VALUES = [
    'UNKNOWN',
    'DIFFERENT_ACTION_KEY',
    'DIFFERENT_DEPS',
    'DIFFERENT_ENVIRONMENT',
    'DIFFERENT_FILES',
    'CORRUPTED_CACHE_ENTRY',
    'NOT_CACHED',

] as const;
export type MissDetailTuple = typeof ALL_STATUS_VALUES;
export type MissDetailEnum = MissDetailTuple[number];

interface Props {
    status: MissDetailEnum;
}

const STATUS_TAGS: { [key in MissDetailEnum]: React.ReactNode } = {
    UNKNOWN: (
        <Tag icon={<QuestionCircleFilled />} color="default" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Unknown</span>
        </Tag>
    ),
    DIFFERENT_ACTION_KEY: (
        <Tag icon={<KeyOutlined />} color="blue" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Different Action Key</span>
        </Tag>
    ),
    DIFFERENT_DEPS: (
        <Tag icon={<StopOutlined />} color="pink" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Different Dependency</span>
        </Tag>
    ),
    DIFFERENT_ENVIRONMENT: (
        <Tag icon={<InfoCircleOutlined />} color="purple" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Different Environment</span>
        </Tag>
    ),
    DIFFERENT_FILES: (
        <Tag icon={<FileUnknownOutlined />} color="cyan" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Different Files</span>
        </Tag>
    ),
    CORRUPTED_CACHE_ENTRY: (
        <Tag icon={<CloseCircleFilled />} color="orange" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Corrupted Cache Entry</span>
        </Tag>
    ),
    NOT_CACHED: (
        <Tag icon={<MinusCircleFilled />} color="red" className={themeStyles.tag}>
            <span className={themeStyles.tagContent}>Not Cached</span>
        </Tag>
    ),

};

const MissDetailTag: React.FC<Props> = ({ status }) => {
    const resultTag = STATUS_TAGS[status] || STATUS_TAGS.UNKNOWN;
    return <>{resultTag}</>;
};

export default MissDetailTag;
