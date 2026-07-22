import { CodeFilled } from "@ant-design/icons";
import { Table } from "antd";
import type React from "react";
import {
  OperationFilterSelector,
  type OperationStatus,
} from "@/components/OperationFilterSelector";
import OperationsInvocationFilter from "@/components/OperationsInvocationFilter";
import { PortalCard } from "@/components/PortalCard";
import type { OperationState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { OperationsFilterParams } from "@/routes/operations.index";
import themeStyles from "@/theme/theme.module.css";
import getColumns from "./Columns";

export const PAGE_SIZE = 1000;

interface Props {
  operations: OperationState[];
  filter: OperationsFilterParams;
  statusFilter: OperationStatus;
  onStatusFilterChange: (value: OperationStatus) => void;
}

export const OperationsPage: React.FC<Props> = ({
  operations,
  filter,
  statusFilter,
  onStatusFilterChange,
}) => {
  return (
    <PortalCard
      icon={<CodeFilled />}
      titleBits={[<span key="title">Operations</span>]}
    >
      <div
        style={{ display: "flex", justifyContent: "flex-end", marginBottom: 8 }}
      >
        <OperationFilterSelector
          value={statusFilter}
          onChange={onStatusFilterChange}
        />
      </div>
      <OperationsInvocationFilter filter={filter} />
      <Table
        dataSource={operations}
        columns={getColumns()}
        pagination={{ pageSize: PAGE_SIZE, hideOnSinglePage: true }}
        size="small"
        rowClassName={() => themeStyles.compactTable}
        rowKey={(item) => item.name}
        locale={{
          emptyText: "No active operations found (that you have access to).",
        }}
      />
    </PortalCard>
  );
};
