import { useQuery } from "@apollo/client";
import { Row, Space } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import React from "react";
import {
  CursorTable,
  getNewPaginationVariables,
} from "@/components/CursorTable";
import type { PaginationVariables } from "@/components/CursorTable/types";
import type { TargetWhereInput } from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import PortalAlert from "../../PortalAlert";
import TargetGridRow from "../TargetGridRow";
import { columns, type TargetGridRowType } from "./Columns";
import { GET_TARGETS_LIST } from "./graphql";

const TargetGrid: React.FC = () => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] =
    React.useState<TargetWhereInput>({});

  const { data, loading, error } = useQuery(GET_TARGETS_LIST, {
    variables: {
      ...paginationVariables,
      where: filterVariables,
    },
  });

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: TargetWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "target-kind":
            newFilters.push({ targetKindContainsFold: value[0] as string });
            break;
          case "label":
            newFilters.push({ labelContainsFold: value[0] as string });
            break;
        }
      }
    });
    setFilterVariables({ and: newFilters });
  };

  if (error) {
    return (
      <PortalAlert
        className="error"
        message="There was a problem communicating w/the backend server."
      />
    );
  }

  const rowData = parseGraphqlEdgeList(data?.findTargets);

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      <Row>
        <CursorTable<TargetGridRowType>
          columns={columns}
          loading={loading}
          size="small"
          rowKey="id"
          expandable={{
            rowExpandable: (_) => true,
            expandedRowRender: (record) => (
              <TargetGridRow
                instanceName={record.instanceName.name}
                label={record.label}
                aspect={record.aspect}
                targetKind={record.targetKind}
                numberOfElements={40}
                direction={"newToOld"}
              />
            ),
          }}
          onChange={(_pagination, filters, _sorter, _extra) =>
            onFilterChange(filters)
          }
          dataSource={rowData}
          pagination={{
            position: "bottom",
            justify: "end",
            size: "middle",
          }}
          pageInfo={{
            startCursor: data?.findTargets.pageInfo.startCursor || "",
            endCursor: data?.findTargets.pageInfo.endCursor || "",
            hasNextPage: data?.findTargets.pageInfo.hasNextPage || false,
            hasPreviousPage:
              data?.findTargets.pageInfo.hasPreviousPage || false,
          }}
          paginationVariables={paginationVariables}
          setPaginationVariables={setPaginationVariables}
        />
      </Row>
    </Space>
  );
};

export default TargetGrid;
