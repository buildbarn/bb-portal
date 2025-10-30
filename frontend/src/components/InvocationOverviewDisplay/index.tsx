import React, { RefAttributes } from 'react';
import { CardProps, Descriptions, Space } from 'antd';
import { JSX } from 'react/jsx-runtime';
import IntrinsicAttributes = JSX.IntrinsicAttributes;
import BuildStepResultTag, { BuildStepResultEnum } from '../BuildStepResultTag';
import PortalDuration from '../PortalDuration';

interface Props {
    command: string,
    cpu: string,
    user: string,
    status: string,
    invocationId: string,
    instanceName: string | undefined,
    configuration: string
    startedAt: string
    endedAt: string
    hostname: string
    isCiWorker: boolean,
    stepLabel: string,
    numFetches: number,
    bazelVersion: string,
}

export const InvocationOverviewDisplay: React.FC<Props> = ({
    command,
    cpu,
    user,
    status,
    invocationId,
    instanceName,
    configuration,
    startedAt,
    endedAt,
    hostname,
    isCiWorker,
    stepLabel,
    numFetches,
    bazelVersion,
}) => {
    return (
        <Space>
            <Descriptions column={1} bordered >
                <Descriptions.Item label="Status">
                    <BuildStepResultTag key="result" result={status as BuildStepResultEnum} />
                </Descriptions.Item>
                <Descriptions.Item label="Invocation Id">
                    {invocationId}
                </Descriptions.Item>
                {instanceName != undefined && instanceName !== "" &&
                    <Descriptions.Item label="Instance name">
                        {instanceName}
                    </Descriptions.Item>
                }
                <Descriptions.Item label="Duration">
                    <PortalDuration key="duration" from={startedAt} to={endedAt} includeIcon includePopover />
                </Descriptions.Item>
                {user != "" &&
                    <Descriptions.Item label="User">
                        {user}
                    </Descriptions.Item>
                }
                {command != "" &&
                    <Descriptions.Item label="Command">
                        <code>
                            {command}
                        </code>
                    </Descriptions.Item>
                }
                {cpu != "" &&
                    <Descriptions.Item label="CPU">
                        {cpu}
                    </Descriptions.Item>
                }
                {configuration != "" &&
                    <Descriptions.Item label="Configuration">
                        {configuration}
                    </Descriptions.Item>
                }
                {hostname != "" &&
                    <Descriptions.Item label="Hostname">
                        {hostname}
                    </Descriptions.Item>
                }
                {numFetches != 0 &&
                    <Descriptions.Item label="Number of Fetches">
                        {numFetches}
                    </Descriptions.Item>
                }
                {isCiWorker &&
                    <Descriptions.Item label="CI Worker">True</Descriptions.Item>
                }
                {stepLabel != "" &&
                    <Descriptions.Item label="CI Step Label">{stepLabel}</Descriptions.Item>
                }
                                {bazelVersion != "" &&
                    <Descriptions.Item label="Bazel version">{bazelVersion}</Descriptions.Item>
                }
            </Descriptions>
        </Space>
    );
};

export default InvocationOverviewDisplay;
