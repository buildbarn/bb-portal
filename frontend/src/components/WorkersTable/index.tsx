import { Flex, Row, Space, Table } from "antd";
import type React from "react";
import type {
  PaginationInfo,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import themeStyles from "@/theme/theme.module.css";
import WorkersTablePageSelector from "../WorkersTablePageSelector";
import WorkersTableTypeSelector from "../WorkersTableTypeSelector";
import getColumns from "./Columns";
import { Route, type WorkerListStatus } from "@/routes/scheduler.worker";

type Props = {
  workerStatusFilter: WorkerListStatus;
  data: WorkerState[] | undefined;
  paginationInfo: PaginationInfo | undefined;
  isLoading: boolean;
  pageSize: number;
};

const WorkersTable: React.FC<Props> = ({
  workerStatusFilter,
  data,
  paginationInfo,
  isLoading,
  pageSize,
}) => {
  const navigate = Route.useNavigate()

  const handleFilterChange = (value: WorkerListStatus) => {
    navigate({
      search: (prev) => ({
        ...prev,
        workerStatusFilter: value,
        cursor: undefined
      })
    })
  };

  const goToNextPage = () => {
    const lastId = data?.at(-1)?.id;
    if (!lastId) return;

    navigate({
      search: (prev) => ({
        ...prev,
        cursor: lastId, 
      }),
    })
  };

  const goToFirstPage = () => {
    navigate({
      search: (prev) => ({
        ...prev,
        cursor: undefined, 
      }),
    });
  };

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      <Row>
        <Flex style={{ width: "100%" }} justify="space-between" wrap>
          <WorkersTableTypeSelector
            workerStatusFilter={workerStatusFilter}
            setWorkerStatusFilter={handleFilterChange}
          />
          {paginationInfo && (
            <WorkersTablePageSelector
              paginationInfo={paginationInfo}
              goToFirstPage={goToFirstPage}
              goToNextPage={goToNextPage}
              pageSize={pageSize}
            />
          )}
        </Flex>
      </Row>
      <Row>
        <Table
          columns={getColumns()}
          loading={isLoading}
          bordered={true}
          style={{ width: "100%" }}
          dataSource={data}
          size="small"
          rowClassName={() => themeStyles.compactTable}
          pagination={false}
          rowKey={(item) => Object.entries(item.id).sort().join(",")}
          locale={{
            emptyText:
              "No workers matching the given criteria found (that you have access to).",
          }}
        />
      </Row>
    </Space>
  );
};

export default WorkersTable;
