"use client";

import { DeploymentUnitOutlined } from "@ant-design/icons";
import { useQuery } from "@apollo/client";
import { Space, Typography } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import dayjs from "dayjs";
import Link from "next/link";
import type React from "react";
import { useState } from "react";
import { validate as uuidValidate } from "uuid";
import styles from "@/components/AppBar/index.module.css";
import CollapsableInvocationTimeline from "@/components/CollapsableInvocationTimeline";
import PortalCard from "@/components/PortalCard";
import {
  BazelInvocationOrderField,
  type BazelInvocationWhereInput,
  type FindBuildByUuidQuery,
  type GetBuildInvocationFragment,
  OrderDirection,
  type SourceControlWhereInput,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeListWithFragment } from "@/utils/parseGraphqlEdgeList";
import { CursorTable, getNewPaginationVariables } from "../CursorTable";
import type { PaginationVariables } from "../CursorTable/types";
import { applyInvocationResultTagFilter } from "../InvocationResultTag/filters";
import PortalAlert from "../PortalAlert";
import { columns } from "./Columns";
import {
  GET_BUILD_BY_UUID_QUERY,
  GET_BUILD_INVOCATION_FRAGMENT,
} from "./graphql";

type BuildType = NonNullable<FindBuildByUuidQuery["getBuild"]>;

const getTitleBits = (build: BuildType | undefined): React.ReactNode[] => {
  if (!build) {
    return [];
  }
  return [
    <span key="build">
      Build ID:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight}>
        {build.buildUUID}
      </Typography.Text>
    </span>,
    <span key="copy-icon" className={styles.copyIcon}>
      <Typography.Text copyable={{ text: build.buildUUID }} />
    </span>,
    <span key="build-url">
      Build URL:{" "}
      <Link
        href={build.buildURL}
        target="_blank"
        className={styles.normalWeight}
      >
        <Typography.Text type="secondary">{build.buildURL}</Typography.Text>
      </Link>
    </span>,
  ];
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

export const BuildDetails: React.FC<Props> = ({ buildUUID }) => {
  const [paginationVariables, setPaginationVariables] =
    useState<PaginationVariables>(getNewPaginationVariables());

  const [filterVariables, setFilterVariables] =
    useState<BazelInvocationWhereInput>({});

  const { data, loading, error } = useQuery(GET_BUILD_BY_UUID_QUERY, {
    variables: {
      ...paginationVariables,
      where: filterVariables,
      orderBy: {
        direction: OrderDirection.Desc,
        field: BazelInvocationOrderField.StartedAt,
      },
      buildUUID: buildUUID,
    },
  });

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    let newFilters: BazelInvocationWhereInput[] = [];
    const sourceControllFilters: SourceControlWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "workflow": {
            sourceControllFilters.push({
              workflowContainsFold: value[0] as string,
            });
            break;
          }
          case "job": {
            sourceControllFilters.push({ jobContainsFold: value[0] as string });
            break;
          }
          case "action": {
            sourceControllFilters.push({
              actionContainsFold: value[0] as string,
            });
            break;
          }
          case "invocationID": {
            const invocationID = value[0] as string;
            if (uuidValidate(invocationID)) {
              newFilters.push({ invocationID: invocationID });
            }
            break;
          }
          case "status": {
            newFilters = newFilters.concat(
              applyInvocationResultTagFilter(value),
            );
            break;
          }
        }
      }
    });
    if (sourceControllFilters.length > 0) {
      newFilters.push({
        hasSourceControlWith: sourceControllFilters,
      });
    }
    setFilterVariables({ and: newFilters });
  };

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

  const build = data?.getBuild ?? undefined;
  const invocations = parseGraphqlEdgeListWithFragment(
    GET_BUILD_INVOCATION_FRAGMENT,
    data?.getBuild?.invocations,
  );

  return (
    <PortalCard
      bordered={false}
      icon={<DeploymentUnitOutlined />}
      titleBits={getTitleBits(build)}
      extraBits={getExtraBits(build)}
    >
      <Space direction="vertical" style={{ width: "100%" }}>
        {invocations.length > 1 && (
          <CollapsableInvocationTimeline invocations={invocations.toReversed()} />
        )}
        <CursorTable<GetBuildInvocationFragment>
          columns={columns}
          loading={loading}
          size="small"
          rowKey="id"
          onChange={(_pagination, filters, _sorter, _extra) =>
            onFilterChange(filters)
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
