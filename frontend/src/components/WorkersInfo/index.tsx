import { Descriptions, type DescriptionsProps, type MenuProps } from "antd";
import type React from "react";
import MenuText from "@/components/MenuText";
import type { SizeClassQueueName } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import PropertyTagList from "../PropertyTagList";
import type { PropertyTagListEntry } from "../PropertyTagList/types";

export interface DrainMenus {
  instanceNamePrefix: MenuProps["items"];
  platformProperty: (entry: PropertyTagListEntry) => MenuProps["items"];
}

interface Props {
  sizeClassQueueName: SizeClassQueueName;
  size?: DescriptionsProps["size"];
  drainMenus?: DrainMenus;
}

const WorkersInfo: React.FC<Props> = ({
  sizeClassQueueName,
  size,
  drainMenus,
}) => {
  const instanceNamePrefix =
    sizeClassQueueName.platformQueueName?.instanceNamePrefix;
  return (
    <Descriptions column={1} bordered size={size}>
      <Descriptions.Item label="Instance name prefix">
        {drainMenus?.instanceNamePrefix ? (
          <MenuText items={drainMenus.instanceNamePrefix}>
            {instanceNamePrefix}
          </MenuText>
        ) : (
          instanceNamePrefix
        )}
      </Descriptions.Item>
      <Descriptions.Item label="Platform properties">
        <PropertyTagList
          propertyList={
            sizeClassQueueName.platformQueueName?.platform?.properties
          }
          tagMenu={drainMenus?.platformProperty}
        />
      </Descriptions.Item>
      <Descriptions.Item label="Size class">
        {sizeClassQueueName.sizeClass}
      </Descriptions.Item>
    </Descriptions>
  );
};

export default WorkersInfo;
