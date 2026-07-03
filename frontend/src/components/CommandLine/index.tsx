import { CodeOutlined, InfoCircleOutlined } from "@ant-design/icons";
import { Card, Divider, Tooltip } from "antd";
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
  const filteredOptions = canonicalCommandLine?.options.filter(
    (x) => x.option !== "config",
  );

  return (
    <PortalCard
      icon={<CodeOutlined />}
      titleBits={["Command Line", rawCommand]}
    >
      {parsedOptions && (
        <Card
          size="small"
          title={
            <strong>
              <Tooltip title="The expanded command line options before normalization">
                Parsed Options <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          styles={{
            header: {
              background: "transparent",
            },
            body: {
              padding: 8,
            },
          }}
        >
          {parsedOptions.options.map((option, index) => (
            <div
              key={option}
              style={{
                paddingTop: "4px",
                paddingLeft: "5px",
              }}
            >
              {option}
              {index !== parsedOptions.options.length - 1 && (
                <Divider size="small" />
              )}
            </div>
          ))}
        </Card>
      )}
      {filteredOptions && (
        <Card
          size="small"
          title={
            <strong>
              <Tooltip title="The expanded command line options used by bazel flags after normalization">
                Normalized Options <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          styles={{
            header: {
              background: "transparent",
            },
            body: {
              padding: 8,
            },
          }}
        >
          {filteredOptions.map((item, index) => (
            <div
              key={`${item.option}-${item.value}`}
              style={{
                paddingTop: "4px",
                paddingLeft: "5px",
              }}
            >
              --{item.option}={item.value}
              {index !== filteredOptions.length - 1 && <Divider size="small" />}
            </div>
          ))}
        </Card>
      )}
      {canonicalCommandLine?.startupOptions && (
        <Card
          size="small"
          title={
            <strong>
              <Tooltip title="The startup options for the bazel server process">
                Startup Options <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          styles={{
            header: {
              background: "transparent",
            },
            body: {
              padding: 8,
            },
          }}
        >
          {canonicalCommandLine.startupOptions.map((item, index) => (
            <div
              key={`${item.option}-${item.value}`}
              style={{
                paddingTop: "8px",
                paddingLeft: "5px",
              }}
            >
              --{item.option}={item.value}
              {index !== canonicalCommandLine.startupOptions.length - 1 && (
                <Divider size="small" />
              )}
            </div>
          ))}
        </Card>
      )}
      {environmentVariables && (
        <Card
          size="small"
          title={
            <strong>
              <Tooltip title="The environment variables that the Bazel process was started with. The environment variables are censored to avoid revealing sensitive values">
                Environment Variables <InfoCircleOutlined />
              </Tooltip>
            </strong>
          }
          styles={{
            header: {
              background: "transparent",
            },
            body: {
              padding: 8,
            },
          }}
        >
          {Object.entries(environmentVariables).map((item, index) => (
            <div
              key={`${item[0]}-${item[1]}`}
              style={{
                paddingTop: "8px",
                paddingLeft: "5px",
              }}
            >
              {item[0]}={item[1]}
              {index !== Object.entries(environmentVariables).length - 1 && (
                <Divider size="small" />
              )}
            </div>
          ))}
        </Card>
      )}
    </PortalCard>
  );
};

export default CommandLineDisplay;
