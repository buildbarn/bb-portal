import { DeploymentUnitOutlined, InfoCircleOutlined } from "@ant-design/icons";
import { useQuery } from "@apollo/client/react";
import { Flex, Popover, Space, Tag, Typography } from "antd";
import dayjs from "dayjs";
import type React from "react";
import { useMemo, useState } from "react";
import styles from "@/components/AppBar/index.module.css";
import CollapsableInvocationTimeline from "@/components/CollapsableInvocationTimeline";
import Content from "@/components/Content";
import { OptionalLinkWrapper } from "@/components/OptionalLinkWrapper";
import PortalCard from "@/components/PortalCard";
import {
  BazelInvocationOrderField,
  type BazelInvocationWhereInput,
  type FindBuildByUuidQuery,
  type GetBuildInvocationFragment,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { applyTableFilters } from "@/utils/applyColumnFilters";
import { env } from "@/utils/env";
import {
  parseGraphqlEdgeList,
  parseGraphqlEdgeListWithFragment,
} from "@/utils/parseGraphqlEdgeList";
import { shouldPollInvocation } from "@/utils/shouldPollInvocation";
import { CursorTable, getNewPaginationVariables } from "../../CursorTable";
import type { PaginationVariables } from "../../CursorTable/types";
import PortalAlert from "../../PortalAlert";
import { getColumns } from "./Columns";
import {
  GET_BUILD_BY_UUID_QUERY,
  GET_BUILD_INVOCATION_FRAGMENT,
} from "./graphql";

type BuildType = NonNullable<FindBuildByUuidQuery["getBuild"]>;

const getTitleBits = (build: BuildType | undefined): React.ReactNode[] => {
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

const getExtraBits = (build: BuildType | undefined): React.ReactNode[] => {
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
  buildUUID: string;
}

export const BuildDetailsPage: React.FC<Props> = (params) => {
  return <Content content={<BuildDetails {...params} />} />;
};

const BuildDetails: React.FC<Props> = ({ buildUUID }) => {
  const [paginationVariables, setPaginationVariables] =
    useState<PaginationVariables>(getNewPaginationVariables());

  const [filterVariables, setFilterVariables] = useState<
    BazelInvocationWhereInput[]
  >([]);

  const { data, loading, error } = useQuery(GET_BUILD_BY_UUID_QUERY, {
    variables: {
      ...paginationVariables,
      where: { and: filterVariables },
      orderBy: {
        direction: OrderDirection.Desc,
        field: BazelInvocationOrderField.StartedAt,
      },
      buildUUID: buildUUID,
    },
  });

  const tableColumns = useMemo(getColumns, []);

  const build = data?.getBuild ?? undefined;
  const tags = parseGraphqlEdgeList(build?.tags);
  const invocations = parseGraphqlEdgeListWithFragment(
    GET_BUILD_INVOCATION_FRAGMENT,
    data?.getBuild?.invocations,
  );
  const inProgressInvocations = invocations
    .filter((inv) => shouldPollInvocation(inv))
    .map((inv) => inv.id);

  // Refetch any ongoing invocations periodically. The result of the query is
  // unused, but in the background Apollo updates the result of the original
  // query based on the IDs of the response.
  useQuery(GET_BUILD_BY_UUID_QUERY, {
    variables: {
      where: {
        idIn: inProgressInvocations,
      },
      buildUUID: buildUUID,
    },
    skip: inProgressInvocations.length === 0,
    pollInterval: 5000,
  });

  if (error) {
    return (
      <PortalCard icon={<DeploymentUnitOutlined />} titleBits={["Build"]}>
        <PortalAlert
          showIcon
          type="error"
          message="Error fetching build details"
          description={
            <>
              {error?.message ? (
                <Typography.Paragraph>{error.message}</Typography.Paragraph>
              ) : (
                <Typography.Paragraph>
                  Unknown error occurred while fetching data from the server.
                </Typography.Paragraph>
              )}
              <Typography.Paragraph>
                It could be that the build is too old and has been removed, or
                that you don&quot;t have access to this build.
              </Typography.Paragraph>
            </>
          }
        />
      </PortalCard>
    );
  }

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
        <CursorTable<GetBuildInvocationFragment>
          columns={tableColumns}
          loading={loading}
          size="small"
          rowKey="id"
          onChange={(_pagination, filters, _sorter, _extra) =>
            applyTableFilters(tableColumns, filters, setFilterVariables)
          }
          dataSource={invocations}
          pagination={{
            position: "bottom",
            justify: "end",
            size: "middle",
          }}
          pageInfo={{
            startCursor: data?.getBuild?.invocations.pageInfo.startCursor || "",
            endCursor: data?.getBuild?.invocations.pageInfo.endCursor || "",
            hasNextPage:
              data?.getBuild?.invocations.pageInfo.hasNextPage || false,
            hasPreviousPage:
              data?.getBuild?.invocations.pageInfo.hasPreviousPage || false,
          }}
          paginationVariables={paginationVariables}
          setPaginationVariables={setPaginationVariables}
        />
      </Space>
    </PortalCard>
  );
};
