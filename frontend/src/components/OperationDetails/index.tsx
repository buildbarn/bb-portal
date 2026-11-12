import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Space, Spin } from "antd";
import { Code } from "@/lib/grpc-client/google/rpc/code";
import { useScheduler } from "@/schedulers/context";
import {
  killPermissionOptions,
  operationDetailsOptions,
} from "@/schedulers/queries";
import ExecuteResponseDisplay from "../ExecutionResponseDisplay";
import OperationStateDisplay from "../OperationStateDisplay";
import PortalAlert from "../PortalAlert";

interface Props {
  operationID: string;
}

const OperationDetails: React.FC<Props> = ({ operationID }) => {
  const { scheduler, client } = useScheduler();
  const queryClient = useQueryClient();

  const killOperationMutation = useMutation({
    mutationFn: () => {
      return client.killOperations({
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
        queryKey: ["operationDetails", scheduler.id, operationID],
        exact: true,
      });
    },
  });

  const { data, isLoading, isError, error } = useQuery(
    operationDetailsOptions(scheduler.id, operationID),
  );

  const { data: allowedToKillOperation } = useQuery(
    killPermissionOptions(scheduler.id, operationID),
  );

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
        title="Error fetching operation"
        description={errorMessage}
      />
    );
  }

  return (
    <Space orientation="vertical" size="middle" style={{ display: "flex" }}>
      {data?.operation && (
        <OperationStateDisplay
          operation={data.operation}
          schedulerId={scheduler.id}
        />
      )}
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
