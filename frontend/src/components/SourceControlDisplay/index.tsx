import { Descriptions, Space } from "antd";
import type React from "react";
import type { BazelInvocationSourceControlFragment } from "@/graphql/__generated__/graphql";
import { OptionalLinkWrapper } from "../OptionalLinkWrapper";

const SourceControlDisplay: React.FC<{
  sourceControl: BazelInvocationSourceControlFragment[];
}> = ({ sourceControl }) => {
  return (
    <Space size="large" direction="vertical">
      {sourceControl.map((sc) => (
        <Descriptions bordered column={1} key={sc.id}>
          {sc.repo ? (
            <Descriptions.Item label="Repository">
              <OptionalLinkWrapper url={sc.repoURL || undefined}>
                {sc.repo || ""}
              </OptionalLinkWrapper>
            </Descriptions.Item>
          ) : undefined}
          {sc.ref ? (
            <Descriptions.Item label={"Ref"}>
              <OptionalLinkWrapper url={sc.refURL || undefined}>
                {sc.ref || ""}
              </OptionalLinkWrapper>
            </Descriptions.Item>
          ) : undefined}
          {sc.commit ? (
            <Descriptions.Item label="Commit SHA">
              <OptionalLinkWrapper url={sc.commitURL || undefined}>
                {sc.commit || ""}
              </OptionalLinkWrapper>
            </Descriptions.Item>
          ) : undefined}
        </Descriptions>
      ))}
    </Space>
  );
};

export default SourceControlDisplay;
