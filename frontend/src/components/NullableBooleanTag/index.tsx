import React from 'react';
import { CheckCircleFilled, CloseCircleFilled, QuestionCircleFilled, } from '@ant-design/icons';
import { Tag } from 'antd';
import themeStyles from '@/theme/theme.module.css';

interface Props {
    status: boolean | null;
    hideText?: boolean
}
export const BOOL_MAP = [
    "true_tag",
    "false_tag",
    "false_hide_tag",
    "null_tag",
    "true_hide_tag"
] as const;

export type BoolTuple = typeof BOOL_MAP;
export type NilBoolEnum = BoolTuple[number]

const BOOL_TAGS: { [key in NilBoolEnum]: React.ReactNode } = {
    false_tag: (
        <Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag}>No</Tag>
    ),
    false_hide_tag: (<Tag icon={<CloseCircleFilled />} color="red" className={themeStyles.tag} />),
    true_tag: (
        <Tag icon={<CheckCircleFilled />} color="green" className={themeStyles.tag}>Yes</Tag>
    ),
    true_hide_tag: (<Tag icon={<CheckCircleFilled />} color="green" className={themeStyles.tag} />),
    null_tag: (
        <Tag icon={<QuestionCircleFilled />} color="orange" className={themeStyles.tag}>?</Tag>
    ),
};

const NullBooleanTag: React.FC<Props> = ({ status, hideText }) => {
    if (hideText == null) {
        hideText = false
    }
    var status_string: NilBoolEnum = "null_tag";
    if (status == true && hideText == false) {
        status_string = "true_tag"
    }
    if (status == true && hideText == true) {
        status_string = "true_hide_tag"
    }
    if (status == false && hideText == false) {
        status_string = "false_tag"
    }
    if (status == false && hideText == true) {
        status_string = "false_hide_tag"
    }
    const resultTag = BOOL_TAGS[status_string] || BOOL_TAGS.null_tag;
    return <>{resultTag}</>;
};

export default NullBooleanTag;
