import { useGrpcClients } from "@/context/GrpcClientsContext";
import { useQuery } from "@tanstack/react-query";
import { Row, Space, Statistic, Typography } from "antd";
import type React from "react";
import PlatformQueuesTable from "../PlatformQueuesTable";
import PortalAlert from "../PortalAlert";

const SchedulerGrid: React.FC = () => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isError, error } = useQuery({
    queryKey: ["listOperations"],
    queryFn: buildQueueStateClient.listOperations.bind(window, {}),
  });

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Row>
        {isError ? (
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
        ) : (
          <Statistic
            title="Total number of operations"
            value={data?.paginationInfo?.totalEntries}
          />
        )}
      </Row>
      <Row>
        <PlatformQueuesTable />
      </Row>
    </Space>
  );
};

export default SchedulerGrid;
