import { DeploymentUnitOutlined, InfoCircleOutlined } from "@ant-design/icons";
import { Flex, Popover, Space, Tag, Typography } from "antd";
import dayjs from "dayjs";
import type React from "react";
import { useMemo } from "react";
import styles from "@/components/AppBar/index.module.css";
import CollapsableInvocationTimeline from "@/components/CollapsableInvocationTimeline";
import { OptionalLinkWrapper } from "@/components/OptionalLinkWrapper";
import { PageCursorTable } from "@/components/PageCursorTable";
import type { GetPaginationUpdateLinkType } from "@/components/PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "@/components/PageCursorTable/utils";
import PortalCard from "@/components/PortalCard";
import type {
  BazelInvocationWhereInput,
  GetBuildFragment,
  GetBuildInvocationFragment,
  PageInfo,
} from "@/graphql/__generated__/graphql";
import { env } from "@/utils/env";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import { getColumns } from "./Columns";

const getTitleBits = (
  build: GetBuildFragment | undefined,
): React.ReactNode[] => {
  if (!build) {
    return [];
  }
  const titleBits = [];
  titleBits.push(
    <Space direction="horizontal" size={0}>
      <Typography.Title level={5}>{`Build ID:`}</Typography.Title>
      <Typography.Text
        copyable
        type="secondary"
        className={styles.normalWeight}
      >
        {build.buildUUID}
      </Typography.Text>
    </Space>,
  );

  const tags = parseGraphqlEdgeList(build.tags);
  const additionalColumns = env.additionalBuildColumns;
  for (const column of additionalColumns) {
    const valueTags = tags.filter((tag) => tag.key === column.valueKey);
    const urlTags = tags.filter((tag) => tag.key === column.urlKey);
    const urlTag = urlTags.length === 1 ? urlTags[0] : undefined;

    titleBits.push(
      <Space direction="horizontal" size={0}>
        <Typography.Title level={5}>{`${column.title}:`}</Typography.Title>
        <Space direction="horizontal">
          <OptionalLinkWrapper url={urlTag?.value}>
            <Typography.Text type="secondary" className={styles.normalWeight}>
              {valueTags.map((tag) => tag.value).join(", ")}
            </Typography.Text>
          </OptionalLinkWrapper>
          {urlTags.length > 1 && (
            <Popover
              title="This field has multiple urls:"
              content={urlTags.map((tag) => (
                <a key={tag.id} href={tag.value}>
                  {tag.value}
                </a>
              ))}
            >
              <InfoCircleOutlined />
            </Popover>
          )}
        </Space>
      </Space>,
    );
  }

  return titleBits;
};

const getExtraBits = (
  build: GetBuildFragment | undefined,
): React.ReactNode[] => {
  if (!build) {
    return [];
  }
  return [
    <Typography.Text key="timestamp" code ellipsis className={styles.startedAt}>
      {dayjs(build.timestamp).format("YYYY-MM-DD hh:mm:ss A")}
    </Typography.Text>,
  ];
};

interface Props {
  invocations: GetBuildInvocationFragment[];
  build: GetBuildFragment;
  pageSize: number;
  onFilterChange: (where: BazelInvocationWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

export const BuildDetailsPage: React.FC<Props> = ({
  invocations,
  build,
  pageSize,
  onFilterChange,
  getPaginationUpdateLink,
  pageInfo,
}) => {
  const tableColumns = useMemo(getColumns, []);
  const tags = parseGraphqlEdgeList(build?.tags);

  return (
    <PortalCard
      variant="outlined"
      icon={<DeploymentUnitOutlined />}
      titleBits={getTitleBits(build)}
      extraBits={getExtraBits(build)}
    >
      <Space direction="vertical" style={{ width: "100%" }}>
        {tags && tags.length > 0 && (
          <Flex gap="4px 0" wrap>
            {tags?.map((tag) => (
              <Tag
                color="blue"
                key={`${tag.key}:${tag.value}`}
                style={{ fontWeight: "bold" }}
              >
                {tag.key}: {tag.value}
              </Tag>
            ))}
          </Flex>
        )}
        {invocations.length > 1 && (
          <CollapsableInvocationTimeline
            invocations={invocations.toReversed()}
          />
        )}
        <PageCursorTable<GetBuildInvocationFragment>
          columns={tableColumns}
          size="small"
          rowKey="id"
          onChange={(_pagination, filters, _sorter, _extra) =>
            onFilterChange(tableFiltersToGraphqlWhere(tableColumns, filters))
          }
          dataSource={invocations}
          pageInfo={pageInfo}
          pageSize={pageSize}
          getPaginationUpdateLink={getPaginationUpdateLink}
        />
      </Space>
    </PortalCard>
  );
};
