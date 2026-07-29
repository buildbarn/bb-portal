import { Descriptions } from "antd";
import type React from "react";
import type { BazelInvocationTagsFragment } from "@/graphql/__generated__/graphql";

const linkRegex =
  /^(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)$/;

interface Props {
  tags: BazelInvocationTagsFragment[];
}

export const InvocationTagTab: React.FC<Props> = ({ tags }) => {
  return (
    <Descriptions bordered column={1} style={{ width: "max-content" }}>
      {tags?.map((tag) => (
        <Descriptions.Item key={tag.key} label={tag.key}>
          {linkRegex.test(tag.value) ? (
            <a href={tag.value}>{tag.value}</a>
          ) : (
            tag.value
          )}
        </Descriptions.Item>
      ))}
    </Descriptions>
  );
};
