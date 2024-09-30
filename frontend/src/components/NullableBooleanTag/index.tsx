import React from 'react';
import { CheckCircleFilled, CloseCircleFilled, QuestionCircleFilled, } from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

interface Props {
    status: boolean | null;
}
export const BOOL_MAP = [
    "true_tag",
    "false_tag",
    "null_tag"
] as const;

export type BoolTuple = typeof BOOL_MAP;
export type NilBoolEnum = BoolTuple[number]

const BOOL_TAGS: { [key in NilBoolEnum]: React.ReactNode } = {
    false_tag: (
        <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>No</Tag>

    ),
    true_tag: (
        <Tag icon={<CheckCircleFilled />} color="green" className={themeStyles.tag}>Yes</Tag>

    ),
    null_tag: (
        <Tag icon={<QuestionCircleFilled />} color="orange" className={themeStyles.tag}>?</Tag>

    ),
};

const NullBooleanTag: React.FC<Props> = ({ status }) => {
    var status_string: NilBoolEnum = "null_tag";
    if (status == true) {
        status_string = "true_tag"
    }
    if (status == false) {
        status_string = "false_tag"
    }
    const resultTag = BOOL_TAGS[status_string] || BOOL_TAGS.null_tag;
    return <>{resultTag}</>;
};

export default NullBooleanTag;
