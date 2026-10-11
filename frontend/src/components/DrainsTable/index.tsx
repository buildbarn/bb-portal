import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Space, Table, Typography } from "antd";
import type React from "react";
import AddDrainForm from "@/components/AddDrainForm";
import PortalAlert from "@/components/PortalAlert";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import themeStyles from "@/theme/theme.module.css";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";
import getColumns from "./Columns";
import { drainKey, modifyDrainsPermissionUrl } from "./utils";

interface Props {
  sizeClassQueueName: SizeClassQueueNameSearch;
}

const DrainsTable: React.FC<Props> = ({ sizeClassQueueName }) => {
  const queryClient = useQueryClient();
  const instanceNamePrefix =
    sizeClassQueueName.platformQueueName.instanceNamePrefix;

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["listDrains", sizeClassQueueName],
    queryFn: () => buildQueueStateClient.listDrains({ sizeClassQueueName }),
  });

  const { data: allowedToModifyDrains } = useQuery({
    queryKey: ["modifyDrainsButtonState", instanceNamePrefix],
    queryFn: async (): Promise<boolean> => {
      const response = await fetch(
        modifyDrainsPermissionUrl(instanceNamePrefix),
      );
      return (await response.json()).allowed;
    },
  });
  const canModify = allowedToModifyDrains === true;

  const onSuccess = () => {
    queryClient.invalidateQueries({
      queryKey: ["listDrains", sizeClassQueueName],
      exact: true,
    });
    queryClient.invalidateQueries({ queryKey: ["listPlatformQueues"] });
  };

  const addDrainMutation = useMutation({
    mutationFn: (workerIdPattern: Record<string, string>) =>
      buildQueueStateClient.addDrain({ sizeClassQueueName, workerIdPattern }),
    onSuccess,
  });

  const removeDrainMutation = useMutation({
    mutationFn: (workerIdPattern: Record<string, string>) =>
      buildQueueStateClient.removeDrain({
        sizeClassQueueName,
        workerIdPattern,
      }),
    onSuccess,
  });

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        title="Error fetching drains"
        description={
          error.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  return (
    <Space orientation="vertical" size="middle" style={{ width: "100%" }}>
      {canModify && (
        <AddDrainForm
          onSubmit={(workerIdPattern, onAdded) =>
            addDrainMutation.mutate(workerIdPattern, { onSuccess: onAdded })
          }
          loading={addDrainMutation.isPending}
        />
      )}
      {allowedToModifyDrains === false && (
        <Typography.Text type="secondary">
          You are not allowed to add or remove drains for this platform queue.
        </Typography.Text>
      )}
      <Table
        columns={getColumns(
          canModify,
          (workerIdPattern) => removeDrainMutation.mutate(workerIdPattern),
          removeDrainMutation.isPending
            ? removeDrainMutation.variables
            : undefined,
        )}
        loading={isLoading}
        bordered={true}
        style={{ width: "100%" }}
        dataSource={data?.drains}
        size="small"
        rowClassName={() => themeStyles.compactTable}
        pagination={false}
        rowKey={(item) => drainKey(item.workerIdPattern)}
        locale={{
          emptyText: "No drains found (that you have access to).",
        }}
      />
      {addDrainMutation.isError && (
        <PortalAlert
          showIcon
          type="error"
          title="Error adding drain"
          description={
            addDrainMutation.error.message ||
            "Unknown error occurred while fetching data from the server."
          }
        />
      )}
      {removeDrainMutation.isError && (
        <PortalAlert
          showIcon
          type="error"
          title="Error removing drain"
          description={
            removeDrainMutation.error.message ||
            "Unknown error occurred while fetching data from the server."
          }
        />
      )}
    </Space>
  );
};

export default DrainsTable;
