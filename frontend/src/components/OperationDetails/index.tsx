import { useGrpcClients } from "@/context/GrpcClientsContext";
import themeStyles from "@/theme/theme.module.css";
import { CloseCircleOutlined, CodeOutlined } from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Space, Typography } from "antd";
import { GrpcErrorCodes } from "../../utils/grpcErrorCodes";
import ExecuteResponseDisplay from "../ExecutionResponseDisplay";
import OperationStateDisplay from "../OperationStateDisplay";

interface Props {
  label: string;
}

const OperationDetails: React.FC<Props> = ({ label }) => {
  const { buildQueueStateClient } = useGrpcClients();
  const queryClient = useQueryClient();

  const killOperationMutation = useMutation({
    mutationFn: () => {
      return buildQueueStateClient.killOperations({
        filter: {
          operationName: label,
        },
        status: {
          code: GrpcErrorCodes.UNAVAILABLE,
          message: "Operation was killed through the web interface",
        },
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["operationDetails", label],
        exact: true,
      });
    },
  });

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["operationDetails", label],
    queryFn: buildQueueStateClient.getOperation.bind(window, {
      operationName: label,
    }),
    staleTime: Number.POSITIVE_INFINITY,
    refetchOnMount: "always",
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
      {data?.operation?.completed ? (
        <ExecuteResponseDisplay executeResponse={data.operation.completed} />
      ) : (
        <Button
          variant="filled"
          color="red"
          danger
          onClick={() => killOperationMutation.mutate()}
        >
          Kill operation
        </Button>
      )}
    </Space>
  );
};

export default OperationDetails;
