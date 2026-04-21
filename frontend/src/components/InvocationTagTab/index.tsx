import { TagsOutlined } from "@ant-design/icons";
import { Descriptions } from "antd";
import type React from "react";
import type { InvocationTag } from "@/graphql/__generated__/graphql";
import PortalCard from "../PortalCard";

const linkRegex =
  /^(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)$/;

interface Props {
  tags: Omit<InvocationTag, "bazelInvocation">[] | undefined;
}

export const InvocationTagTab: React.FC<Props> = ({ tags }) => {
  return (
    <PortalCard
      type="inner"
      icon={<TagsOutlined />}
      titleBits={[<span key="title">Tags</span>]}
      style={{ width: "100%" }}
    >
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
    </PortalCard>
  );
};
