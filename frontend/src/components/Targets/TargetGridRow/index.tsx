import { useQuery } from "@apollo/client";
import { Space } from "antd";
import type React from "react";
import {
  InvocationTargetOrderField,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import TargetGridBtn from "../TargetGridBtn";
import { GET_INVOCATION_TARGETS_FOR_TARGET } from "./graphql";

interface Props {
  instanceName: string;
  label: string;
  aspect: string;
  targetKind: string;
  numberOfElements: number;
  direction: "oldToNew" | "newToOld";
}

const TargetGridRow: React.FC<Props> = ({
  instanceName,
  label,
  aspect,
  targetKind,
  numberOfElements,
  direction: reverseOrder,
}) => {
  var { data } = useQuery(GET_INVOCATION_TARGETS_FOR_TARGET, {
    variables: {
      instanceName,
      label,
      aspect,
      targetKind,
      first: numberOfElements,
      orderBy: {
        field: InvocationTargetOrderField.StartedAt,
        direction: OrderDirection.Desc,
      },
    },
    fetchPolicy: "cache-and-network",
  });

  const rowData = parseGraphqlEdgeList(data?.getTarget?.invocationTargets);
  if (reverseOrder === "oldToNew") {
    rowData.reverse();
  }

  return (
    <Space size={0} style={{ paddingLeft: "40px" }}>
      {rowData.map((item) => (
        <TargetGridBtn
          key={`target-grid-btn-${item.id}`}
          invocationId={item.bazelInvocation.invocationID}
          status={item.success ?? null}
        />
      ))}
    </Space>
  );
};

export default TargetGridRow;
