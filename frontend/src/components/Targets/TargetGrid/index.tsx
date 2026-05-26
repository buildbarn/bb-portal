import { Row, Space } from "antd";
import type React from "react";
import { PageCursorTable } from "@/components/PageCursorTable";
import type { GetPaginationUpdateLinkType } from "@/components/PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "@/components/PageCursorTable/utils";
import type {
  PageInfo,
  TargetListDetailsFragment,
  TargetWhereInput,
} from "@/graphql/__generated__/graphql";
import TargetGridRow from "../TargetGridRow";
import { columns } from "./Columns";

interface Props {
  targets: TargetListDetailsFragment[];
  onFilterChange: (where: TargetWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageSize: number;
  pageInfo: PageInfo;
}

const TargetGrid: React.FC<Props> = ({
  targets,
  onFilterChange,
  getPaginationUpdateLink,
  pageSize,
  pageInfo,
}) => {
  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Row>
        <PageCursorTable<TargetListDetailsFragment>
          columns={columns}
          size="small"
          rowKey="id"
          expandable={{
            rowExpandable: (_) => true,
            expandedRowRender: (record) => (
              <TargetGridRow
                targetId={record.id}
                numberOfElements={40}
                direction={"newToOld"}
              />
            ),
          }}
          onChange={(_pagination, filters, _sorter, _extra) =>
            onFilterChange(tableFiltersToGraphqlWhere(columns, filters))
          }
          dataSource={targets}
          pageInfo={pageInfo}
          getPaginationUpdateLink={getPaginationUpdateLink}
          pageSize={pageSize}
        />
      </Row>
    </Space>
  );
};

export default TargetGrid;
