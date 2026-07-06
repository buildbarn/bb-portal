import { Descriptions } from "antd";
import type React from "react";
import type { BazelInvocationOverviewFragment } from "@/graphql/__generated__/graphql";
import { commandLineDataToString } from "@/utils/commandLineDataToString";
import { InvocationResultTag } from "../InvocationResultTag";
import PortalDuration from "../PortalDuration";

interface Props {
  invocation: BazelInvocationOverviewFragment;
}

export const InvocationOverviewDisplay: React.FC<Props> = ({ invocation }) => {
  const {
    invocationID,
    startedAt,
    endedAt,
    exitCodeName,
    configurations,
    instanceName,
    connectionMetadata,
    originalCommandLine,
    numFetches,
    hostname,
    bazelVersion,
  } = invocation;

  const command = commandLineDataToString(originalCommandLine);

  // TODO: Determine how to best display multiple configurations
  const cpu = Array.from(
    new Set(
      configurations
        ?.map((config) => config.cpu)
        ?.filter((cpu) => cpu && cpu !== ""),
    ),
  )
    .sort()
    .join(", ");
  const mnemonics = Array.from(
    new Set(
      configurations
        ?.map((config) => config.mnemonic)
        ?.filter((mnemonic) => mnemonic && mnemonic !== ""),
    ),
  )
    .sort()
    .join(", ");

  return (
    <Descriptions column={1} bordered style={{ width: "max-content" }}>
      <Descriptions.Item label="Status">
        <InvocationResultTag
          key="result"
          exitCodeName={exitCodeName}
          timeSinceLastConnectionMillis={
            connectionMetadata?.timeSinceLastConnectionMillis
          }
        />
      </Descriptions.Item>
      <Descriptions.Item label="Invocation Id">
        {invocationID}
      </Descriptions.Item>
      {instanceName.name !== "" && (
        <Descriptions.Item label="Instance name">
          {instanceName.name}
        </Descriptions.Item>
      )}
      <Descriptions.Item label="Duration">
        <PortalDuration
          key="duration"
          from={startedAt || undefined}
          to={
            endedAt
              ? endedAt
              : connectionMetadata?.connectionLastOpenAt || undefined
          }
          includeIcon
          formatConfig={{ smallestUnit: "s" }}
        />
      </Descriptions.Item>
      {command !== "" && (
        <Descriptions.Item label="Command">
          <code>{command}</code>
        </Descriptions.Item>
      )}
      {cpu !== "" && <Descriptions.Item label="CPU">{cpu}</Descriptions.Item>}
      {mnemonics !== "" && (
        <Descriptions.Item label="Configuration mnemonics">
          {mnemonics}
        </Descriptions.Item>
      )}
      {hostname !== "" && (
        <Descriptions.Item label="Hostname">{hostname}</Descriptions.Item>
      )}
      {numFetches !== 0 && (
        <Descriptions.Item label="Number of Fetches">
          {numFetches}
        </Descriptions.Item>
      )}
      {bazelVersion !== "" && (
        <Descriptions.Item label="Bazel version">
          {bazelVersion}
        </Descriptions.Item>
      )}
    </Descriptions>
  );
};
