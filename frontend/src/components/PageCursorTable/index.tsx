import { LeftOutlined, RightOutlined } from "@ant-design/icons";
import { useNavigate } from "@tanstack/react-router";
import { Col, Row, Select, Space, Table, type TableProps } from "antd";
import type { AnyObject } from "antd/es/_util/type";
import type React from "react";
import { LinkButton } from "../LinkButton";
import type { GetPaginationUpdateLinkType, PageInfo } from "./types";

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
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
}

export function PageCursorTable<RecordType = AnyObject>({
  children,
  pageInfo,
  pageSize,
  getPaginationUpdateLink,
  ...antdTableProps
}: Props<RecordType>) {
  const navigate = useNavigate();
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
            <LinkButton
              {...getPaginationUpdateLink({
                pagination: {
                  after: undefined,
                  first: undefined,
                  before: pageInfo.startCursor,
                  last: pageSize,
                },
              })}
              size="middle"
              disabled={!pageInfo.hasPreviousPage}
            >
              <LeftOutlined />
            </LinkButton>
            <LinkButton
              {...getPaginationUpdateLink({
                pagination: {
                  after: pageInfo.endCursor,
                  first: pageSize,
                  before: undefined,
                  last: undefined,
                },
              })}
              size="middle"
              disabled={!pageInfo.hasNextPage}
            >
              <RightOutlined />
            </LinkButton>
            <Select
              defaultValue={pageSize}
              style={{ width: 120 }}
              options={getPageSizeOptions()}
              onChange={(value) =>
                navigate(
                  getPaginationUpdateLink({
                    pageSize: value,
                    pagination: {
                      after: undefined,
                      first: value,
                      before: undefined,
                      last: undefined,
                    },
                  }),
                )
              }
            />
          </Space>
        </Col>
      </Row>
    </Space>
  );
}
