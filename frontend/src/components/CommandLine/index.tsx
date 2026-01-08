import { CodeOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { Empty, List, Row, Tooltip, Typography } from 'antd';
import React from 'react';
import PortalCard from '../PortalCard';

// TODO: find a way to apply these interfaces automatically to the
// output of graphql while remaining a scalar with regard to the graphql
// api.
interface CommandLineData {
  executable: string
  command: string
  options: CommandLineOptions[]
  startupOptions: CommandLineOptions[]
  residual: string[]
}

interface CommandLineOptions {
  option: string
  value: string
}

interface ParsedOptions {
  explicitOptions: string[]
  options: string[]
}

interface Props {
  rawCommand: string | null
  canonicalCommandLine: CommandLineData | null
  parsedOptions: ParsedOptions
}

const CommandLineDisplay: React.FC<Props> = ({rawCommand, canonicalCommandLine, parsedOptions}) => {
  if (!canonicalCommandLine) {
    return <PortalCard icon={<CodeOutlined />} titleBits={["Command Line", rawCommand]}>
      <Empty description="No information about the command line available..." />
    </PortalCard>;
  }
  
  const opts = canonicalCommandLine.options.filter(x => x.option !== 'config')

  return (
    <PortalCard icon={<CodeOutlined />} titleBits={["Command Line", rawCommand]}>
      { !parsedOptions ? null :
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
        renderItem={x => <List.Item>{x}</List.Item>}
        />
      }
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
        dataSource={opts}
        renderItem={(item) => <List.Item>--{item.option}={item.value}</List.Item>}
      />
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
        renderItem={(item) => <List.Item>--{item.option}={item.value}</List.Item>}
      />
    </PortalCard>
  );
}

export default CommandLineDisplay;
