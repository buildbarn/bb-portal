import React from 'react';
import {
    CheckCircleFilled,
    CloseCircleFilled,
    InfoCircleFilled,
    MinusCircleFilled,
    QuestionCircleFilled,
    StopOutlined,
} from '@ant-design/icons';
import { Button, Tag } from 'antd';
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
    invocationId: string;
}


const TestGridBtn: React.FC<Props> = ({ status, invocationId }) => {

    const ICON_BTNS: { [key in TestStatusEnum]: React.ReactNode } = {
        NO_STATUS: (
            <Button icon={<QuestionCircleFilled />} className={themeStyles.colorDisabled} />

        ),
        PASSED: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<CheckCircleFilled />} color="success" className={themeStyles.colorSuccess} />

        ),
        FLAKY: (
            <Button icon={<InfoCircleFilled />} color="warning" className={themeStyles.colorAborted} />

        ),
        FAILED: (
            <Button icon={<CloseCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        TIMEOUT: (
            <Button icon={<MinusCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        INCOMPLETE: (
            <Button icon={<StopOutlined />} color="secondary" className={themeStyles.colorAborted} />

        ),
        REMOTE_FAILURE: (
            <Button icon={<CloseCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        FAILED_TO_BUILD: (
            <Button icon={<QuestionCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        TOOL_HALTED_BEFORE_TESTING: (
            <Button icon={<QuestionCircleFilled />} color="secondary" className={themeStyles.colorDisabled} />
        ),
    };

    const resultTag = ICON_BTNS[status] || ICON_BTNS.NO_STATUS;
    return <>{resultTag}</>;
};

export default TestGridBtn;
