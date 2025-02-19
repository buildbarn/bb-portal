import { Flex, Tag } from "antd";
import type React from "react";
import type { PropertyTagListEntry } from "./types";

interface Props {
  propertyList: PropertyTagListEntry[] | undefined;
}

const PropertyTagList: React.FC<Props> = ({ propertyList }) => {
  return (
    <Flex gap="4px 0" wrap>
      {propertyList?.map((entry) => (
        <Tag
          color="blue"
          key={`${entry.name}:${entry.value}`}
          style={{ fontWeight: "bold" }}
        >
          {entry.name}: {entry.value}
        </Tag>
      ))}
    </Flex>
  );
};

export default PropertyTagList;
