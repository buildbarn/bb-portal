import { Table } from "antd";
import type React from "react";
import themeStyles from "@/theme/theme.module.css";
import { useOperatorMode } from "@/utils/operatorMode";
import PortalAlert from "../PortalAlert";
import { usePlatformQueues, usePlatformQueuesBulkDrain } from "./bulkDrain";
import getColumns from "./Columns";

const PlatformQueuesTable: React.FC = () => {
  const [operatorMode] = useOperatorMode();
  const { data, isLoading, isError, error } = usePlatformQueues();
  const { bulk, contextHolder } = usePlatformQueuesBulkDrain(
    data,
    operatorMode,
  );

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        title="Error fetching platform queues"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <>
      {contextHolder}
      <Table
        columns={getColumns(operatorMode, bulk)}
        loading={isLoading}
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
          emptyText: "No platform queues found (that you have access to).",
        }}
      />
    </>
  );
};

export default PlatformQueuesTable;
