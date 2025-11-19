import { useQuery } from "@tanstack/react-query";
import { Skeleton, Space, Statistic } from "antd";
import Link from "next/link";
import type React from "react";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import PortalAlert from "../PortalAlert";

export const SchedulerStatistics: React.FC = () => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isError, error, isPending } = useQuery({
    queryKey: ["listOperations"],
    queryFn: buildQueueStateClient.listOperations.bind(null, {}),
  });

  if (isPending) {
    return (
      <Space direction="vertical" size={9.5}>
        <Skeleton.Node active style={{ width: 180, height: 16 }} />
        <Skeleton.Node active style={{ width: 180, height: 32 }} />
      </Space>
    );
  }

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching scheduler statistics"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Statistic
      title="Total number of operations"
      value={data.paginationInfo?.totalEntries}
      valueRender={(value) => <Link href="/operations">{value}</Link>}
    />
  );
};
