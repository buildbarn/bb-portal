import { useGrpcClients } from "@/context/GrpcClientsContext";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Space, Spin } from "antd";
import { env } from "next-runtime-env";
import { Code } from "@/lib/grpc-client/google/rpc/code";
import ExecuteResponseDisplay from "../ExecutionResponseDisplay";
import OperationStateDisplay from "../OperationStateDisplay";
import PortalAlert from "../PortalAlert";

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
          code: Code.UNAVAILABLE,
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
        `${env("NEXT_PUBLIC_BES_BACKEND_URL") || ""}/api/v1/checkPermissions/killOperation/${operationID}`,
      );
      return (await response.json()).allowed;
    },
  });

  if (isLoading) {
    return <Spin />;
  }

  if (isError) {
    let errorMessage: string;
    if (
      error.message ===
      "/buildbarn.buildqueuestate.BuildQueueState/GetOperation NOT_FOUND: Operation was not found"
    ) {
      errorMessage =
        "The operation was not found. Operations are automatically removed 60 seconds after they have been completed. It may have already completed and been removed, or you may not have access to view it.";
    } else {
      errorMessage =
        error.message ||
        "Unknown error occurred while fetching data from the server.";
    }
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching operation"
        description={errorMessage}
      />
    );
  }

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
