"use client";

import React, { useState } from "react";
import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";
import { Space, Table, TableColumnsType, Typography } from "antd";
import {
  DeploymentUnitOutlined,
  FilterOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import TargetDetails from "@/components/Targets/TargetDetails";
import {
  FindBuildByUuidQuery,
  FindBuildByUuidQueryVariables,
} from "@/graphql/__generated__/graphql";
import { FIND_BUILD_BY_UUID_QUERY } from "./index.graphql";
import { useQuery } from "@apollo/client";
import PortalAlert from "@/components/PortalAlert";
import Link from "next/link";
import SearchWidget, { SearchFilterIcon } from "@/components/SearchWidgets";
import styles from "@/components/AppBar/index.module.css";
import dayjs from "dayjs";
import BuildStepResultTag, {
  BuildStepResultEnum,
} from "@/components/BuildStepResultTag";
import PortalDuration from "@/components/PortalDuration";

interface StatusProps {
  buildUUID: string;
}

interface BuildGridRowDataType {
  key: React.Key;
  user: string;
  invocationId: string;
  startedAt: string;
  endedAt: string;
  status: string;
  workflow: string;
  action: string;
  job: string;
}

export default function Page({ params }: { params: { buildUUID: string } }) {
  const [variables, setVariables] = useState<FindBuildByUuidQueryVariables>({});

  const {
    loading: loading,
    data: responseData,
    error: error,
    previousData: previousData,
  } = useQuery(FIND_BUILD_BY_UUID_QUERY, {
    variables: { uuid: params.buildUUID },
    pollInterval: 60000,
  });

  const data = loading ? previousData : responseData;
  var result: BuildGridRowDataType[] = [];

  if (error) {
    console.error(error);
    return (
      <PortalAlert
        className="error"
        message="There was a problem communicating w/the backend server."
      />
    );
  } else {
    data?.getBuild?.invocations?.map((row) => {
      result.push({
        key: row.id ?? "",
        user: row.userLdap ?? "",
        invocationId: row.invocationID ?? "",
        startedAt: row.startedAt ?? "",
        endedAt: row.endedAt ?? "",
        status: row.state?.exitCode?.name ?? "",
        workflow: row.sourceControl?.workflow ?? "",
        action: row.sourceControl?.action ?? "",
        job: row.sourceControl?.job ?? "",
      });
    });
  }

  const titleBits: React.ReactNode[] = [
    <span key="build">
      Build ID:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight}>
        {data?.getBuild?.buildUUID}
      </Typography.Text>
    </span>,
  ];
  titleBits.push(
    <span className={styles.copyIcon}>
      <Typography.Text
        copyable={{ text: data?.getBuild?.buildUUID ?? "Copy" }}
      />
    </span>
  );
  titleBits.push(
    <span key="build-url">
      Build URL:{" "}
      <Link
        href={data?.getBuild?.buildURL ?? ""}
        target="_blank"
        className={styles.normalWeight}
      >
        <Typography.Text type="secondary">
          {data?.getBuild?.buildURL}
        </Typography.Text>{" "}
      </Link>
    </span>
  );
  const extraBits: React.ReactNode[] = [
    <Typography.Text code ellipsis className={styles.startedAt}>
      {dayjs(data?.getBuild?.timestamp).format("YYYY-MM-DD hh:mm:ss A")}
    </Typography.Text>,
  ];

  const workflow_filters: string[] = Array.from(
    new Set(
      data?.getBuild?.invocations?.map(
        (x) => x.sourceControl?.workflow ?? ""
      ) ?? []
    )
  );

  const job_filters: string[] = Array.from(
    new Set(
      data?.getBuild?.invocations?.map((x) => x.sourceControl?.job ?? "") ?? []
    )
  );

  const action_filters: string[] = Array.from(
    new Set(
      data?.getBuild?.invocations?.map((x) => x.sourceControl?.action ?? "") ??
        []
    )
  );

  const columns: TableColumnsType<BuildGridRowDataType> = [
    {
      title: "Workflow",
      dataIndex: "workflow",
      filterSearch: true,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.workflow.includes(value.toString()) ? true : false,
      sorter: (a, b) => a.workflow.localeCompare(b.workflow),
      filters: workflow_filters.map((x) => ({ text: x, value: x })),
    },
    {
      title: "Job",
      dataIndex: "job",
      filterSearch: true,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.job.includes(value.toString()) ? true : false,
      sorter: (a, b) => a.job.localeCompare(b.job),
      filters: job_filters.map((x) => ({ text: x, value: x })),
    },
    {
      title: "Action",
      dataIndex: "action",
      filterSearch: true,
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.action.includes(value.toString()) ? true : false,
      sorter: (a, b) => a.action.localeCompare(b.action),
      filters: action_filters.map((x) => ({ text: x, value: x })),
    },

    {
      title: "User",
      dataIndex: "user",
    },
    {
      title: "Invocation ID",
      dataIndex: "invocationId",
      filterSearch: true,
      filterDropdown: (filterProps) => (
        <SearchWidget placeholder="Target Pattern..." {...filterProps} />
      ),
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) =>
        record.invocationId.includes(value.toString()) ? true : false,
      render: (_, record) => (
        <Space>
          <span className={styles.copyIcon}>
            <Typography.Text
              copyable={{ text: record.invocationId ?? "Copy" }}
            />
          </span>
          <Link href={"/bazel-invocations/" + record.invocationId}>
            {record.invocationId}
          </Link>
        </Space>
      ),
    },
    {
      title: "Duration",
      dataIndex: "startedAt",
      render: (_, record) => (
        <PortalDuration
          key="duration"
          from={record.startedAt}
          to={record.endedAt}
          includeIcon
          includePopover
        />
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      filterSearch: true,
      render: (_, record) => (
        <BuildStepResultTag
          key="result"
          result={record.status as BuildStepResultEnum}
        />
      ),
      filterIcon: (filtered) => (
        <SearchFilterIcon icon={<FilterOutlined />} filtered={filtered} />
      ),
      onFilter: (value, record) => record.status == value,
      filters: [
        {
          text: "Succeeded",
          value: "SUCCESS",
        },
        {
          text: "Unstable",
          value: "UNSTABLE",
        },
        {
          text: "Parsing Failed",
          value: "PARSING_FAILURE",
        },
        {
          text: "Build Failed",
          value: "BUILD_FAILURE",
        },
        {
          text: "Tests Failed",
          value: "TESTS_FAILED",
        },
        {
          text: "Not Built",
          value: "NOT_BUILT",
        },
        {
          text: "Aborted",
          value: "ABORTED",
        },
        {
          text: "Interrupted",
          value: "INTERRUPTED",
        },
        {
          text: "Status Unknown",
          value: "UNKNOWN",
        },
      ],
    },
  ];

  return (
    <Content
      content={
        <PortalCard
          bordered={false}
          icon={<DeploymentUnitOutlined />}
          titleBits={titleBits}
          extraBits={extraBits}
        >
          <Space direction="vertical" style={{ width: "100%" }}>
            <Table columns={columns} dataSource={result} pagination={false} />
          </Space>
        </PortalCard>
      }
    />
  );
}
