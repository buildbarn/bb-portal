import { LeftOutlined, RightOutlined } from "@ant-design/icons";
import { Button, Col, Row, Select, Space, Table, type TableProps } from "antd";
import type { AnyObject } from "antd/es/_util/type";
import type { SizeType } from "antd/es/config-provider/SizeContext";
import type React from "react";
import type { PageInfo, PaginationVariables } from "./types";

export const PAGE_SIZE_OPTIONS = [50, 100, 200, 500, 1000];

const PAGE_SIZE_DEFAULT = PAGE_SIZE_OPTIONS[0];

export const getNewPaginationVariables = (): PaginationVariables => {
  return {
    pageSize: PAGE_SIZE_DEFAULT,
    first: PAGE_SIZE_DEFAULT,
  };
};

const getPageSizeOptions = () => {
  return PAGE_SIZE_OPTIONS.map((size) => ({
    value: size,
    label: `${size} / page`,
  }));
};

interface CursorPaginationProps {
  size: SizeType;
  pageInfo: PageInfo;
  paginationVariables: PaginationVariables;
  setPaginationVariables: React.Dispatch<
    React.SetStateAction<PaginationVariables>
  >;
}

const CursorPagination: React.FC<CursorPaginationProps> = ({
  size,
  pageInfo,
  paginationVariables,
  setPaginationVariables,
}) => {
  return (
    <Space direction="horizontal">
      <Button
        size={size}
        disabled={!pageInfo.hasPreviousPage}
        onClick={() =>
          setPaginationVariables({
            pageSize: paginationVariables.pageSize,
            before: pageInfo.startCursor,
            last: paginationVariables.pageSize,
            after: undefined,
            first: undefined,
          })
        }
      >
        <LeftOutlined />
      </Button>

      <Button
        size={size}
        disabled={!pageInfo.hasNextPage}
        onClick={() =>
          setPaginationVariables({
            pageSize: paginationVariables.pageSize,
            after: pageInfo.endCursor,
            first: paginationVariables.pageSize,
            before: undefined,
            last: undefined,
          })
        }
      >
        <RightOutlined />
      </Button>
      <Select
        defaultValue={paginationVariables.pageSize}
        style={{ width: 120 }}
        onChange={(value) =>
          setPaginationVariables({
            pageSize: value,
            first: value,
            after: undefined,
            before: undefined,
            last: undefined,
          })
        }
        options={getPageSizeOptions()}
      />
    </Space>
  );
};

interface CursorTableProps<RecordType = AnyObject>
  extends Omit<TableProps<RecordType>, "pagination"> {
  children?: React.ReactNode;
  pageInfo: PageInfo;
  pagination: {
    position: "top" | "bottom";
    justify: "start" | "center" | "end";
    size: SizeType;
  };
  paginationVariables: PaginationVariables;
  setPaginationVariables: React.Dispatch<
    React.SetStateAction<PaginationVariables>
  >;
}

export function CursorTable<RecordType = AnyObject>({
  children,
  pageInfo,
  pagination,
  paginationVariables,
  setPaginationVariables,
  ...rest
}: CursorTableProps<RecordType>) {
  const paginationButtons = (
    <Row justify={pagination.justify}>
      <Col>
        <CursorPagination
          size={pagination.size}
          pageInfo={pageInfo}
          paginationVariables={paginationVariables}
          setPaginationVariables={setPaginationVariables}
        />
      </Col>
    </Row>
  );

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      {pagination.position === "top" && paginationButtons}

      <Table<RecordType>
        {...(rest as TableProps<RecordType>)}
        pagination={false}
      >
        {children}
      </Table>

      {pagination.position === "bottom" && paginationButtons}
    </Space>
  );
}
