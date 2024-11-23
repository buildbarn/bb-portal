import React, { RefAttributes } from 'react';
import { CardProps, Descriptions, Space } from 'antd';
import { JSX } from 'react/jsx-runtime';
import IntrinsicAttributes = JSX.IntrinsicAttributes;
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
    hostname: string
    isCiWorker: boolean,
    stepLabel: string,
    numFetches: number
}

type OverviewProps = Props & IntrinsicAttributes & CardProps & RefAttributes<HTMLDivElement>;

export const InvocationOverviewDisplay: React.FC<OverviewProps> = ({
    targets,
    command,
    cpu,
    user,
    status,
    invocationId,
    configuration,
    startedAt,
    endedAt,
    hostname,
    isCiWorker,
    stepLabel,
    numFetches, ...props }) => {

        const stepLabelExists = stepLabel != "" && stepLabel != null

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
                    <code>
                        {command}
                    </code>
                </Descriptions.Item>
                <Descriptions.Item label="CPU">
                    {cpu}
                </Descriptions.Item>
                <Descriptions.Item label="Configuration">
                    {configuration}
                </Descriptions.Item>
                <Descriptions.Item label="Hostname">
                    {hostname}
                </Descriptions.Item>
                <Descriptions.Item label="Number of Fetches">
                    {numFetches}
                </Descriptions.Item>
                { isCiWorker &&
                 <Descriptions.Item label="CI Worker">True</Descriptions.Item>
                }
                { stepLabelExists &&
                    <Descriptions.Item label="CI Step Label">{stepLabel}</Descriptions.Item>
                }
            </Descriptions>
        </Space>
    );
};

export default InvocationOverviewDisplay;
