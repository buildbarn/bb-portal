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
            <Button href={"/bazel-invocations/" + invocationId} icon={<CheckCircleFilled />} className={themeStyles.colorSuccess} />

        ),
        FLAKY: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<InfoCircleFilled />} className={themeStyles.colorAborted} />

        ),
        FAILED: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<CloseCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        TIMEOUT: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<MinusCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        INCOMPLETE: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<StopOutlined />} className={themeStyles.colorAborted} />

        ),
        REMOTE_FAILURE: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<CloseCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        FAILED_TO_BUILD: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<QuestionCircleFilled />} color="danger" className={themeStyles.colorFailure} />

        ),
        TOOL_HALTED_BEFORE_TESTING: (
            <Button href={"/bazel-invocations/" + invocationId} icon={<QuestionCircleFilled />} className={themeStyles.colorDisabled} />
        ),
    };

    const resultTag = ICON_BTNS[status] || ICON_BTNS.NO_STATUS;
    return <>{resultTag}</>;
};

export default TestGridBtn;
