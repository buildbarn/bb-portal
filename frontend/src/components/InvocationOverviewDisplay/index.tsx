import React, { RefAttributes } from 'react';
import { Card, CardProps, Descriptions, List, Space } from 'antd';
import { AnsiUp } from 'ansi_up';
import linkifyHtml from 'linkify-html';
import { JSX } from 'react/jsx-runtime';
import styles from './index.module.css';
import PortalAlert from '@/components/PortalAlert';
import IntrinsicAttributes = JSX.IntrinsicAttributes;
import PortalCard from '../PortalCard';
import { ExclamationCircleFilled } from '@ant-design/icons';
import BuildStepResultTag, { BuildStepResultEnum } from '../BuildStepResultTag';
import PortalDuration from '../PortalDuration';

interface Props {

    targets: number,
    command: string,
    cpu: string,
    user: string,
    status: string,
    invocationId: string,
    configuration: string
    startedAt: string
    endedAt: string
    numFetches: number
}

type OverviewProps = Props & IntrinsicAttributes & CardProps & RefAttributes<HTMLDivElement>;

export const InvocationOverviewDisplay: React.FC<OverviewProps> = ({ targets, command, cpu, user, status, invocationId, configuration, startedAt, endedAt, numFetches, ...props }) => {
    return (
        <Space>
            <Descriptions column={1} bordered >
                <Descriptions.Item label="Status">
                    <BuildStepResultTag key="result" result={status as BuildStepResultEnum} />
                </Descriptions.Item>
                <Descriptions.Item label="Invocation Id">
                    {invocationId}
                </Descriptions.Item>
                <Descriptions.Item label="Duration">
                    <PortalDuration key="duration" from={startedAt} to={endedAt} includeIcon includePopover />
                </Descriptions.Item>
                <Descriptions.Item label="User">
                    {user}
                </Descriptions.Item>
                <Descriptions.Item label="Command">
                    {command}
                </Descriptions.Item>
                <Descriptions.Item label="CPU">
                    {cpu}
                </Descriptions.Item>
                <Descriptions.Item label="Configuration">
                    {configuration}
                </Descriptions.Item>
                <Descriptions.Item label="Number of Fetches">
                    {numFetches}
                </Descriptions.Item>
            </Descriptions>
        </Space>
    );
};

export default InvocationOverviewDisplay;
