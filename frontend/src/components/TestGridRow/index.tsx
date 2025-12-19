import { useQuery } from "@apollo/client";
import { Space } from "antd";
import type React from "react";
import {
  type GetTestsForTargetQuery,
  InvocationTargetOrderField,
  OrderDirection,
  TestSummaryOrderField,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import TestGridBtn from "../TestGridBtn";
import type { TestStatusEnum } from "../TestStatusTag";
import { GET_TESTS_FOR_TARGET } from "./graphql";

interface Props {
  instanceName: string;
  label: string;
  aspect: string;
  targetKind: string;
  numberOfElements: number;
  direction: "oldToNew" | "newToOld";
}

const TestGridRow: React.FC<Props> = ({
  instanceName,
  label,
  aspect,
  targetKind,
  numberOfElements,
  direction,
}) => {
  var { data } = useQuery<GetTestsForTargetQuery>(GET_TESTS_FOR_TARGET, {
    variables: {
      first: numberOfElements,
      where: {
        hasInvocationTargetWith: {
          hasTargetWith: {
            hasInstanceNameWith: {
              name: instanceName,
            },
            label: label,
            aspect: aspect,
            targetKind: targetKind,
          },
        },
      },
      orderBy: {
        field: TestSummaryOrderField.FirstStartTime,
        direction:
          direction === "oldToNew" ? OrderDirection.Asc : OrderDirection.Desc,
      },
    },
    fetchPolicy: "cache-and-network",
  });

  const rowData = parseGraphqlEdgeList(data?.findTestSummaries);

  return (
    <Space size={0} style={{ paddingLeft: "40px" }}>
      {rowData.map((item) => (
        <TestGridBtn
          key={"test-grid-btn" + item.id}
          invocationId={item.invocationTarget.bazelInvocation.invocationID}
          status={item.overallStatus as TestStatusEnum}
        />
      ))}
    </Space>
  );
};

export default TestGridRow;
