import { Flex, Row, Space, Table } from "antd";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import type React from "react";
import type {
  PaginationInfo,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import themeStyles from "@/theme/theme.module.css";
import type { ListWorkerFilterType } from "@/types/ListWorkerFilterType";
import WorkersTablePageSelector from "../WorkersTablePageSelector";
import WorkersTableTypeSelector from "../WorkersTableTypeSelector";
import getColumns from "./Columns";

type Props = {
  listWorkerFilterType: ListWorkerFilterType;
  data: WorkerState[] | undefined;
  paginationInfo: PaginationInfo | undefined;
  isLoading: boolean;
  pageSize: number;
};

const WorkersTable: React.FC<Props> = ({
  listWorkerFilterType,
  data,
  paginationInfo,
  isLoading,
  pageSize,
}) => {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const changeUrlQueryValue = (key: string, value: string | undefined) => {
    const params = new URLSearchParams(searchParams.toString());
    if (value === undefined) {
      params.delete(key);
    } else {
      params.set(key, value);
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  const handleFilterChange = (value: ListWorkerFilterType) => {
    changeUrlQueryValue("listWorkerFilterType", value);
  };

  const goToNextPage = () => {
    changeUrlQueryValue("paginationCursor", JSON.stringify(data?.at(-1)?.id));
  };

  const goToFirstPage = () => {
    changeUrlQueryValue("paginationCursor", undefined);
  };

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      <Row>
        <Flex style={{ width: "100%" }} justify="space-between" wrap>
          <WorkersTableTypeSelector
            listWorkerFilterType={listWorkerFilterType}
            setListWorkerFilterType={handleFilterChange}
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
