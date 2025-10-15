import React from 'react';
import { Descriptions, Space } from 'antd';
import PortalDuration from '../PortalDuration';
import { InvocationResultTag } from '../InvocationResultTag';

interface Props {
    command: string,
    cpu: string,
    exitCodeName: string | undefined,
    bepCompleted: boolean,
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
    exitCodeName,
    bepCompleted,
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
                    <InvocationResultTag key="result" exitCodeName={exitCodeName} bepCompleted={bepCompleted} />
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
