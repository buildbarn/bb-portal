import { useGrpcClients } from "@/context/GrpcClientsContext";
import { useQuery } from "@tanstack/react-query";
import { Row, Skeleton, Space, Statistic, Typography } from "antd";
import type React from "react";
import PlatformQueuesTable from "../PlatformQueuesTable";
import PortalAlert from "../PortalAlert";

const SchedulerGrid: React.FC = () => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isError, error, isPending } = useQuery({
    queryKey: ["listOperations"],
    queryFn: buildQueueStateClient.listOperations.bind(window, {}),
  });

  const totalNumberOfOperaionsDisplay = () => {
    if (isError) {
      return (
        <PortalAlert
          className="error"
          message={
            <>
              <Typography.Text>
                There was a problem communicating with the backend server:
              </Typography.Text>
              <pre>{String(error)}</pre>
            </>
          }
        />
      );
    }
    if (isPending) {
      return (
        <Space direction="vertical" size={9.5}>
          <Skeleton.Node active style={{ width: 180, height: 16 }} />
          <Skeleton.Node active style={{ width: 180, height: 32 }} />
        </Space>
      );
    }
    return (
      <Statistic
        title="Total number of operations"
        value={data.paginationInfo?.totalEntries}
      />
    );
  };

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Row>{totalNumberOfOperaionsDisplay()}</Row>
      <Row>
        <PlatformQueuesTable />
      </Row>
    </Space>
  );
};

export default SchedulerGrid;
