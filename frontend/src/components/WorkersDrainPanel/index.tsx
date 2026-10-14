import { PlusOutlined } from "@ant-design/icons";
import { Button, Card, Flex, Input, Space, Tag, Typography } from "antd";
import type React from "react";
import { useState } from "react";
import type { ListedMatch } from "@/components/DrainActions";
import { patternToPropertyList } from "@/components/DrainsTable/utils";

interface Props {
  workerIdPattern: Record<string, string>;
  match: ListedMatch | undefined;
  onSet: (key: string, value: string) => void;
  onRemove: (key: string) => void;
  onClear: () => void;
  onDrain: () => void;
  onDrainQueue: () => void;
}

const WorkersDrainPanel: React.FC<Props> = ({
  workerIdPattern,
  match,
  onSet,
  onRemove,
  onClear,
  onDrain,
  onDrainQueue,
}) => {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const pairs = patternToPropertyList(workerIdPattern);
  const empty = pairs.length === 0;

  const add = () => {
    onSet(key.trim(), value);
    setKey("");
    setValue("");
  };

  return (
    <Card
      size="small"
      title="Drain"
      style={{ width: "100%" }}
      extra={
        <Space size="small">
          <Button size="small" danger onClick={onDrainQueue}>
            Drain this whole queue
          </Button>
          <Button size="small" disabled={empty} onClick={onClear}>
            Clear
          </Button>
          <Button
            size="small"
            type="primary"
            danger
            disabled={empty}
            onClick={onDrain}
          >
            Drain matching workers
          </Button>
        </Space>
      }
    >
      <Flex justify="space-between" align="center" gap="small" wrap>
        <Flex align="center" gap="small" wrap>
          <Typography.Text strong>Pattern:</Typography.Text>
          {empty ? (
            <Typography.Text type="secondary">
              Click worker ID labels below to build a drain pattern.
            </Typography.Text>
          ) : (
            <Flex gap="4px 0" wrap>
              {pairs.map(({ name, value }) => (
                <Tag
                  key={name}
                  color="blue"
                  variant="solid"
                  closable
                  style={{ fontWeight: "bold" }}
                  onClose={(e) => {
                    e.preventDefault();
                    onRemove(name);
                  }}
                >
                  {name}: {value}
                </Tag>
              ))}
            </Flex>
          )}
          {!empty && match && (
            <Typography.Text type="secondary">
              {match.matching} of {match.listed} workers on this page match.
            </Typography.Text>
          )}
        </Flex>
        <Space.Compact size="small">
          <Input
            placeholder="key"
            value={key}
            onChange={(e) => setKey(e.target.value)}
            onPressEnter={() => key.trim() && add()}
            style={{ width: 120 }}
          />
          <Input
            placeholder="value"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            onPressEnter={() => key.trim() && add()}
            style={{ width: 160 }}
          />
          <Button icon={<PlusOutlined />} disabled={!key.trim()} onClick={add}>
            Add
          </Button>
        </Space.Compact>
      </Flex>
      <Typography.Text type="secondary" style={{ fontSize: "smaller" }}>
        To drain every queue sharing a platform property or the instance name
        prefix, click it in the header above.
      </Typography.Text>
    </Card>
  );
};

export default WorkersDrainPanel;
