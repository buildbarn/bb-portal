import { LeftOutlined, RightOutlined } from "@ant-design/icons";
import { Button, Col, Row, Select, Space, Table, type TableProps } from "antd";
import type { AnyObject } from "antd/es/_util/type";
import type React from "react";
import type { OnTablePaginationChange, PageInfo } from "./types";

export const PAGE_SIZE_OPTIONS = [50, 100, 200, 500, 1000];

export const DEFAULT_PAGE_SIZE = PAGE_SIZE_OPTIONS[0];

const getPageSizeOptions = () => {
  return PAGE_SIZE_OPTIONS.map((size) => ({
    value: size,
    label: `${size} / page`,
  }));
};

interface Props<RecordType = AnyObject>
  extends Omit<TableProps<RecordType>, "pagination"> {
  children?: React.ReactNode;
  pageInfo: PageInfo;
  pageSize: number | undefined;
  onPaginationChange: OnTablePaginationChange;
}

export function PageCursorTable<RecordType = AnyObject>({
  children,
  pageInfo,
  pageSize,
  onPaginationChange,
  ...antdTableProps
}: Props<RecordType>) {
  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      <Table<RecordType>
        {...(antdTableProps as TableProps<RecordType>)}
        pagination={false}
      >
        {children}
      </Table>

      <Row justify="end">
        <Col>
          <Space direction="horizontal">
            <Button
              size="middle"
              disabled={!pageInfo.hasPreviousPage}
              // TODO: Find a typesafe way to do this link that supports
              // preloading. A link would be work, but I haven't found a way to
              // make it typesafe with respect to the current route, without
              // hardcoding the route path (which won't work since this is a
              // generic component).
              onClick={() =>
                onPaginationChange({
                  pageSize: pageSize,
                  pagination: {
                    after: undefined,
                    first: undefined,
                    before: pageInfo.startCursor,
                    last: pageSize,
                  },
                })
              }
            >
              <LeftOutlined />
            </Button>
            <Button
              size="middle"
              disabled={!pageInfo.hasNextPage}
              onClick={() =>
                onPaginationChange({
                  pageSize: pageSize,
                  pagination: {
                    after: pageInfo.endCursor,
                    first: pageSize,
                    before: undefined,
                    last: undefined,
                  },
                })
              }
            >
              <RightOutlined />
            </Button>
            <Select
              defaultValue={pageSize}
              style={{ width: 120 }}
              onChange={(value) =>
                onPaginationChange({
                  pageSize: value,
                  pagination: {
                    after: undefined,
                    first: value,
                    before: undefined,
                    last: undefined,
                  },
                })
              }
              options={getPageSizeOptions()}
            />
          </Space>
        </Col>
      </Row>
    </Space>
  );
}
