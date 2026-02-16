"use client";

import { DeploymentUnitOutlined } from "@ant-design/icons";
import { Space, Table, Typography } from "antd";
import dayjs from "dayjs";
import Link from "next/link";
import type React from "react";
import styles from "@/components/AppBar/index.module.css";
import CollapsableInvocationTimeline from "@/components/CollapsableInvocationTimeline";
import PortalCard from "@/components/PortalCard";
import type { FindBuildByUuidQuery } from "@/graphql/__generated__/graphql";
import { getColumns } from "./Columns";

const getTitleBits = (data: FindBuildByUuidQuery): React.ReactNode[] => {
  return [
    <span key="build">
      Build ID:{" "}
      <Typography.Text type="secondary" className={styles.normalWeight}>
        {data?.getBuild?.buildUUID}
      </Typography.Text>
    </span>,
    <span key="copy-icon" className={styles.copyIcon}>
      <Typography.Text
        copyable={{ text: data?.getBuild?.buildUUID ?? "Copy" }}
      />
    </span>,
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
    </span>,
  ];
};

const getExtraBits = (data: FindBuildByUuidQuery): React.ReactNode[] => {
  return [
    <Typography.Text key="timestamp" code ellipsis className={styles.startedAt}>
      {dayjs(data?.getBuild?.timestamp).format("YYYY-MM-DD hh:mm:ss A")}
    </Typography.Text>,
  ];
};

interface Props {
  data: FindBuildByUuidQuery;
}

export const BuildDetails: React.FC<Props> = ({ data }) => {
  const invocations = data?.getBuild?.invocations || [];
  return (
    <PortalCard
      bordered={false}
      icon={<DeploymentUnitOutlined />}
      titleBits={getTitleBits(data)}
      extraBits={getExtraBits(data)}
    >
      <Space direction="vertical" style={{ width: "100%" }}>
        {data?.getBuild?.invocations &&
          data.getBuild.invocations.length > 1 && (
            <CollapsableInvocationTimeline
              invocations={data.getBuild?.invocations}
            />
          )}
        <Table
          columns={getColumns(data)}
          dataSource={invocations}
          pagination={false}
        />
      </Space>
    </PortalCard>
  );
};
