import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { Flex, Row, Space, Table } from "antd";
import type React from "react";
import InvocationChildrenFilterSelector from "@/components/InvocationChildrenFilterSelector";
import PortalAlert from "@/components/PortalAlert";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { ListInvocationChildrenRequest_Filter } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { InvocationChildrenFilter, Route } from "@/routes/scheduler.invocation";
import themeStyles from "@/theme/theme.module.css";
import { isNotFoundGrpcError } from "@/utils/grpcStatus";
import {
  anyToInvocationIdSearch,
  type InvocationNameSearch,
  invocationNameFromSearch,
  rootInvocationNameSearch,
} from "@/utils/invocationName";
import getColumns from "./Columns";

const requestFilters: Record<
  InvocationChildrenFilter,
  ListInvocationChildrenRequest_Filter
> = {
  [InvocationChildrenFilter.QUEUED]:
    ListInvocationChildrenRequest_Filter.QUEUED,
  [InvocationChildrenFilter.ACTIVE]:
    ListInvocationChildrenRequest_Filter.ACTIVE,
  [InvocationChildrenFilter.ALL]: ListInvocationChildrenRequest_Filter.ALL,
};

interface Props {
  invocationName: InvocationNameSearch;
  childrenFilter: InvocationChildrenFilter;
}

const InvocationChildrenTable: React.FC<Props> = ({
  invocationName,
  childrenFilter,
}) => {
  const navigate = Route.useNavigate();

  const { data, isLoading, isPlaceholderData, isError, error } = useQuery({
    queryKey: ["listInvocationChildren", invocationName, childrenFilter],
    queryFn: () =>
      buildQueueStateClient.listInvocationChildren({
        invocationName: invocationNameFromSearch(invocationName),
        filter: requestFilters[childrenFilter],
      }),
    placeholderData: (previousData, _) => previousData,
  });

  if (isError) {
    if (isNotFoundGrpcError(error)) {
      return (
        <Space orientation="vertical">
          <PortalAlert
            showIcon
            type="info"
            title="Invocation no longer exists"
            description="The scheduler removes invocations once they have no queued or executing operations and no idle workers."
          />
          <Link
            to="/scheduler/invocation"
            search={{
              invocationName: rootInvocationNameSearch(
                invocationName.sizeClassQueueName,
              ),
              childrenFilter: InvocationChildrenFilter.ACTIVE,
              cursor: undefined,
            }}
          >
            Go to root invocation
          </Link>
        </Space>
      );
    }
    return (
      <PortalAlert
        showIcon
        type="error"
        title="Error fetching child invocations"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Space orientation="vertical" style={{ width: "100%" }}>
      <Row>
        <Flex style={{ width: "100%" }} justify="space-between" wrap>
          <InvocationChildrenFilterSelector
            childrenFilter={childrenFilter}
            setChildrenFilter={(value) =>
              navigate({
                search: (prev) => ({ ...prev, childrenFilter: value }),
              })
            }
          />
        </Flex>
      </Row>
      <Row>
        <Table
          columns={getColumns(invocationName)}
          loading={isLoading || isPlaceholderData}
          bordered={true}
          style={{ width: "100%" }}
          dataSource={data?.children}
          size="small"
          rowClassName={() => themeStyles.compactTable}
          pagination={false}
          rowKey={(item) => {
            const id = item.id && anyToInvocationIdSearch(item.id);
            return `${id?.typeUrl}:${id?.value}`;
          }}
          locale={{
            emptyText: "No child invocations found (that you have access to).",
          }}
        />
      </Row>
    </Space>
  );
};

export default InvocationChildrenTable;
