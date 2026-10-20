import { CaretDownOutlined } from "@ant-design/icons";
import { Dropdown, Flex, type MenuProps, Tag } from "antd";
import type React from "react";
import type { PropertyTagListEntry } from "./types";

interface Props {
  propertyList: PropertyTagListEntry[] | undefined;
  tagMenu?: (entry: PropertyTagListEntry) => MenuProps["items"];
  onTagClick?: (entry: PropertyTagListEntry) => void;
  isTagSelected?: (entry: PropertyTagListEntry) => boolean;
}

const PropertyTagList: React.FC<Props> = ({
  propertyList,
  tagMenu,
  onTagClick,
  isTagSelected,
}) => {
  return (
    <Flex gap="4px 0" wrap>
      {propertyList?.map((entry) => {
        const key = `${entry.name}:${entry.value}`;
        const style = { fontWeight: "bold", ...entry.style };
        if (onTagClick) {
          return (
            <Tag
              color="blue"
              key={key}
              variant={isTagSelected?.(entry) ? "solid" : undefined}
              style={{ ...style, cursor: "pointer" }}
              onClick={() => onTagClick(entry)}
            >
              {entry.name}: {entry.value}
            </Tag>
          );
        }
        const items = tagMenu?.(entry);
        if (items) {
          return (
            <Dropdown
              key={key}
              trigger={["click", "contextMenu"]}
              menu={{ items }}
            >
              <Tag color="blue" style={{ ...style, cursor: "pointer" }}>
                {entry.name}: {entry.value}{" "}
                <CaretDownOutlined style={{ opacity: 0.6 }} />
              </Tag>
            </Dropdown>
          );
        }
        return (
          <Tag color="blue" key={key} style={style}>
            {entry.name}: {entry.value}
          </Tag>
        );
      })}
    </Flex>
  );
};

export default PropertyTagList;
