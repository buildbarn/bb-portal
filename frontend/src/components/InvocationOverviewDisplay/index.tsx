import React from 'react';
import { Descriptions, Space } from 'antd';
import PortalDuration from '../PortalDuration';
import { InvocationResultTag } from '../InvocationResultTag';
import { Configuration as BazelConfiguration } from '@/graphql/__generated__/graphql';

type Configuration = Pick<BazelConfiguration, 'cpu' | 'mnemonic'>;

interface Props {
    command: string,
    exitCodeName: string | undefined,
    bepCompleted: boolean,
    invocationId: string,
    instanceName: string | undefined,
    configurations: Configuration[] | undefined,
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
    exitCodeName,
    bepCompleted,
    invocationId,
    instanceName,
    configurations,
    startedAt,
    endedAt,
    hostname,
    isCiWorker,
    stepLabel,
    numFetches,
    bazelVersion,
}) => {
    // TODO: Determine how to best display multiple configurations
    const cpu = Array.from(
            new Set(configurations
                ?.map((config) => config.cpu)
                ?.filter((cpu) => cpu && cpu !== "")
            )
        )
        .sort()
        .join(", ");
    let mnemonics = Array.from(
            new Set(configurations
                ?.map((config) => config.mnemonic)
                ?.filter((mnemonic) => mnemonic && mnemonic !== "")
            )
        )
        .sort()
        .join(", ");

    return (
        <Space>
            <Descriptions column={1} bordered >
                <Descriptions.Item label="Status">
                    <InvocationResultTag key="result" exitCodeName={exitCodeName} bepCompleted={bepCompleted} />
                </Descriptions.Item>
                <Descriptions.Item label="Invocation Id">
                    {invocationId}
                </Descriptions.Item>
                {instanceName !== undefined && instanceName !== "" &&
                    <Descriptions.Item label="Instance name">
                        {instanceName}
                    </Descriptions.Item>
                }
                <Descriptions.Item label="Duration">
                    <PortalDuration key="duration" from={startedAt} to={endedAt} includeIcon includePopover />
                </Descriptions.Item>
                {command !== "" &&
                    <Descriptions.Item label="Command">
                        <code>
                            {command}
                        </code>
                    </Descriptions.Item>
                }
                {cpu !== "" &&
                    <Descriptions.Item label="CPU">
                        {cpu}
                    </Descriptions.Item>
                }
                {mnemonics !== "" &&
                    <Descriptions.Item label="Configuration mnemonics">
                        {mnemonics}
                    </Descriptions.Item>
                }
                {hostname !== "" &&
                    <Descriptions.Item label="Hostname">
                        {hostname}
                    </Descriptions.Item>
                }
                {numFetches !== 0 &&
                    <Descriptions.Item label="Number of Fetches">
                        {numFetches}
                    </Descriptions.Item>
                }
                {isCiWorker &&
                    <Descriptions.Item label="CI Worker">True</Descriptions.Item>
                }
                {stepLabel !== "" &&
                    <Descriptions.Item label="CI Step Label">{stepLabel}</Descriptions.Item>
                }
                                {bazelVersion !== "" &&
                    <Descriptions.Item label="Bazel version">{bazelVersion}</Descriptions.Item>
                }
            </Descriptions>
        </Space>
    );
};

export default InvocationOverviewDisplay;
