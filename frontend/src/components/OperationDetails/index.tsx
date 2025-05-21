import { useGrpcClients } from "@/context/GrpcClientsContext";
import themeStyles from "@/theme/theme.module.css";
import { CloseCircleOutlined, CodeOutlined } from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Space, Typography } from "antd";
import { env } from "next-runtime-env";
import { GrpcErrorCodes } from "../../utils/grpcErrorCodes";
import ExecuteResponseDisplay from "../ExecutionResponseDisplay";
import OperationStateDisplay from "../OperationStateDisplay";

interface Props {
  operationID: string;
}

const OperationDetails: React.FC<Props> = ({ operationID }) => {
  const { buildQueueStateClient } = useGrpcClients();
  const queryClient = useQueryClient();

  const killOperationMutation = useMutation({
    mutationFn: () => {
      return buildQueueStateClient.killOperations({
        filter: {
          operationName: operationID,
        },
        status: {
          code: GrpcErrorCodes.UNAVAILABLE,
          message: "Operation was killed through the web interface",
        },
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["operationDetails", operationID],
        exact: true,
      });
    },
  });

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["operationDetails", operationID],
    queryFn: buildQueueStateClient.getOperation.bind(window, {
      operationName: operationID,
    }),
    staleTime: Number.POSITIVE_INFINITY,
    refetchOnMount: "always",
  });

  const { data: allowedToKillOperation } = useQuery({
    queryKey: ["killOperationsButtonState", operationID],
    queryFn: async (): Promise<boolean> => {
      const response = await fetch(
        `${env("NEXT_PUBLIC_BES_BACKEND_URL") || ""}/api/checkPermissions/killOperation/${operationID}`,
      );
      return (await response.json()).allowed;
    },
  });

  if (isError) {
    return (
      <Typography.Text
        disabled
        className={themeStyles.tableEmptyTextTypography}
      >
        <Space direction="vertical">
          <Space>
            <CloseCircleOutlined />
            There was an error fetching the operation details
          </Space>
          <pre>{String(error)}</pre>
        </Space>
      </Typography.Text>
    );
  }

  if (isLoading)
    return (
      <Typography.Text
        disabled
        className={themeStyles.tableEmptyTextTypography}
      >
        <Space>
          <CodeOutlined />
          Loading...
        </Space>
      </Typography.Text>
    );

  return (
    <Space direction="vertical" size="middle" style={{ display: "flex" }}>
      {data?.operation && <OperationStateDisplay operation={data.operation} />}
      {data?.operation?.completed && (
        <ExecuteResponseDisplay executeResponse={data.operation.completed} />
      )}
      {!data?.operation?.completed && (
        <Button
          danger
          disabled={allowedToKillOperation === false}
          onClick={() => killOperationMutation.mutate()}
        >
          Kill operation
        </Button>
      )}
    </Space>
  );
};

export default OperationDetails;
