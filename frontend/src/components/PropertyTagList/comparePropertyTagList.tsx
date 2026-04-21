import { theme } from "antd";
import type React from "react";
import { styleMap } from "../BrowserDirectory/utils";
import PropertyTagList from ".";
import type { PropertyTagListEntry } from "./types";

interface Props {
  propertyList: PropertyTagListEntry[];
  comparePropertyList: PropertyTagListEntry[];
  mergedMode?: boolean; // For comparing actions
}

const { useToken } = theme;

const ComparePropertyTagList: React.FC<Props> = ({
  propertyList,
  comparePropertyList,
  mergedMode,
}) => {
  const { token } = useToken();
  const mergedPropertyList: PropertyTagListEntry[] = mergedMode
    ? [
        ...propertyList.map((entry) => {
          return !comparePropertyList.find(
            (compareEntry) =>
              compareEntry.name === entry.name &&
              compareEntry.value === entry.value,
          )
            ? {
                ...entry,
                style: styleMap("unique", token),
              }
            : entry;
        }),
        ...comparePropertyList
          .filter(
            (entry) =>
              !propertyList?.find(
                (compareEntry) =>
                  compareEntry.name === entry.name &&
                  compareEntry.value === entry.value,
              ),
          )
          .map((entry) => {
            return { ...entry, style: styleMap("missing", token) };
          }),
      ]
    : propertyList.map((entry) => {
        return !comparePropertyList.find(
          (compareEntry) =>
            compareEntry.name === entry.name &&
            compareEntry.value === entry.value,
        )
          ? {
              ...entry,
              style: styleMap("diff_with_borders", token),
            }
          : entry;
      });

  return <PropertyTagList propertyList={mergedPropertyList} />;
};

export default ComparePropertyTagList;
