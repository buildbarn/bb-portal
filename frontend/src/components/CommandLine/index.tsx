import { CodeOutlined, InfoCircleOutlined } from "@ant-design/icons";
import { List, Tooltip } from "antd";
import type React from "react";
import PortalCard from "../PortalCard";

// TODO: find a way to apply these interfaces automatically to the
// output of graphql while remaining a scalar with regard to the graphql
// api.
export interface CommandLineData {
  executable: string;
  command: string;
  options: CommandLineOptions[];
  startupOptions: CommandLineOptions[];
  residual: string[];
}

export interface CommandLineOptions {
  option: string;
  value: string;
}

interface ParsedOptions {
  explicitOptions: string[];
  options: string[];
}

interface Props {
  rawCommand: string | undefined | null;
  canonicalCommandLine: CommandLineData | undefined | null;
  parsedOptions: ParsedOptions | undefined | null;
  environmentVariables: Record<string, string> | undefined | null;
}

const CommandLineDisplay: React.FC<Props> = ({
  rawCommand,
  canonicalCommandLine,
  parsedOptions,
  environmentVariables,
}) => {
  return (
    <PortalCard
      icon={<CodeOutlined />}
      titleBits={["Command Line", rawCommand]}
    >
      {parsedOptions && (
        <List
          bordered
          size="small"
          style={{ width: "100%" }}
          header={
            <strong>
              <Tooltip title="The expanded command line options before normalization">
                Parsed Options <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          dataSource={parsedOptions.options}
          renderItem={(x) => <List.Item>{x}</List.Item>}
        />
      )}
      {canonicalCommandLine?.options && (
        <List
          bordered
          size="small"
          style={{ width: "100%" }}
          header={
            <strong>
              <Tooltip title="The expanded command line options used by bazel flags after normalization">
                Normalized Options <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          dataSource={canonicalCommandLine.options.filter(
            (x) => x.option !== "config",
          )}
          renderItem={(item) => (
            <List.Item>
              --{item.option}={item.value}
            </List.Item>
          )}
        />
      )}
      {canonicalCommandLine?.startupOptions && (
        <List
          bordered
          size="small"
          style={{ width: "100%" }}
          header={
            <strong>
              <Tooltip title="The startup options for the bazel server process">
                Startup Options <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          dataSource={canonicalCommandLine.startupOptions}
          renderItem={(item) => (
            <List.Item>
              --{item.option}={item.value}
            </List.Item>
          )}
        />
      )}
      {environmentVariables && (
        <List
          bordered
          size="small"
          style={{ width: "100%" }}
          header={
            <strong>
              <Tooltip title="The environment variables that the Bazel process was started with. The environment variables are censored to avoid revealing sensitive values.">
                Environment Variables <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          dataSource={Object.entries(environmentVariables).map((item) => ({
            key: item[0],
            value: item[1],
          }))}
          renderItem={(item) => (
            <List.Item>
              {item.key}={item.value}
            </List.Item>
          )}
        />
      )}
    </PortalCard>
  );
};

export default CommandLineDisplay;
