"use client";

import { useQuery } from "@tanstack/react-query";
import { Space, Spin, Typography } from "antd";
import type React from "react";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import { Command } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { fetchCasObjectAndParse } from "@/utils/fetchCasObject";
import BrowserCommandDescription from "../BrowserCommandDescription";
import FilesTable from "../FilesTable";
import { filesTableEntriesFromOutputPath } from "../FilesTable/utils";
import PortalAlert from "../PortalAlert";

interface Params {
  browserPageParams: BrowserPageParams;
}

const BrowserCommandGrid: React.FC<Params> = ({ browserPageParams }) => {
  const { casByteStreamClient } = useGrpcClients();

  const { data, isError, isPending, error } = useQuery({
    queryKey: ["browserCommandGrid", browserPageParams],
    queryFn: () =>
      fetchCasObjectAndParse(
        casByteStreamClient,
        browserPageParams.instanceName,
        browserPageParams.digestFunction,
        browserPageParams.digest,
        Command,
      ),
  });

  if (isPending) {
    return <Spin />;
  }

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching command"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <Typography.Title level={2}>Command</Typography.Title>
      <BrowserCommandDescription
        browserPageParams={browserPageParams}
        command={data}
        commandDigest={browserPageParams.digest}
        showTitle={false}
      />

      <Typography.Title level={2}>Output files</Typography.Title>
      <FilesTable
        entries={data.outputPaths.map((entry) =>
          filesTableEntriesFromOutputPath(entry),
        )}
        isPending={isPending}
      />
    </Space>
  );
};

export default BrowserCommandGrid;
