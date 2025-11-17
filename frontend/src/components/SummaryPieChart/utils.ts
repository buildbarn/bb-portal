import type { GlobalToken } from "antd";

export const themeColor = (token: GlobalToken, index: number): string => {
  const colors: string[] = [
    token["lime-8"],
    token["purple-8"],
    token["red-8"],
    token["cyan-8"],
    token["yellow-8"],
    token["red-6"],
    token["blue-8"],
    token["magenta-8"],
    token["volcano-6"],
    token["green-8"],
    token["volcano-8"],
    token["blue-6"],
    token["orange-8"],
    token["magenta-6"],
    token["green-6"],
    token["orange-6"],
    token["cyan-6"],
  ];

  return colors[index % colors.length];
};
