import { useGrpcClients } from "@/context/GrpcClientsContext";
import themeStyles from "@/theme/theme.module.css";
import { BuildOutlined } from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Space, Table, Typography } from "antd";
import type React from "react";
import getColumns from "./Columns";
import type { PlatformQueueTableState } from "./types";

const PlatformQueuesTable: React.FC = () => {
  const { buildQueueStateClient } = useGrpcClients();

  const { data, isLoading } = useQuery({
    queryKey: ["listPlatformQueues"],
    queryFn: async (): Promise<PlatformQueueTableState[]> => {
      const queues = await buildQueueStateClient.listPlatformQueues({});
      // Convert PlatoformQueueState to PlatformQueueTableState, making it
      // suitable for the table component. This is done by flattening the
      // sizeClassQueues array into a single element array, and adding
      // additional properties to track the number of size classes and if
      // the current size class is the first size class in the queue.
      return queues.platformQueues.flatMap((queue) => {
        return queue.sizeClassQueues.map((sizeClassQueue, index) => {
          return {
            ...queue,
            sizeClassQueues: [sizeClassQueue],
            numberOfSizeClasses: queue.sizeClassQueues.length,
            isFirstSizeClass: index === 0,
          };
        });
      });
    },
  });

  const emptyText = "No platform queues can be found.";

  return (
    <Table
      columns={getColumns()}
      bordered={true}
      style={{ width: "100%" }}
      dataSource={data}
      size="small"
      rowClassName={() => themeStyles.compactTable}
      pagination={false}
      rowKey={(item) =>
        `instanceNamePrefix:${item.name?.instanceNamePrefix}-sizeClass:${item.sizeClassQueues[0].sizeClass}`
      }
      locale={{
        emptyText: isLoading ? (
          <Typography.Text
            disabled
            className={themeStyles.tableEmptyTextTypography}
          >
            <Space>
              <BuildOutlined />
              Loading...
            </Space>
          </Typography.Text>
        ) : (
          <Typography.Text
            disabled
            className={themeStyles.tableEmptyTextTypography}
          >
            <Space>
              <BuildOutlined />
              {emptyText}
            </Space>
          </Typography.Text>
        ),
      }}
    />
  );
};

export default PlatformQueuesTable;
