import type { SizeClassQueueName } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { Descriptions, Space } from "antd";
import type React from "react";
import { JSX } from "react/jsx-runtime";
import PropertyTagList from "../PropertyTagList";
import IntrinsicAttributes = JSX.IntrinsicAttributes;

interface Props {
  sizeClassQueueName: SizeClassQueueName;
}

const WorkersInfo: React.FC<Props> = ({ sizeClassQueueName }) => {
  return (
    <Space>
      <Descriptions column={1} bordered>
        <Descriptions.Item label="Instance name prefix">
          {sizeClassQueueName.platformQueueName?.instanceNamePrefix}
        </Descriptions.Item>
        <Descriptions.Item label="Platform properties">
          <PropertyTagList
            propertyList={
              sizeClassQueueName.platformQueueName?.platform?.properties
            }
          />
        </Descriptions.Item>
        <Descriptions.Item label="Size class">
          {sizeClassQueueName.sizeClass}
        </Descriptions.Item>
      </Descriptions>
    </Space>
  );
};

export default WorkersInfo;
