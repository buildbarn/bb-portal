import { theme as antTheme, Flex, Row, Space, Table } from "antd";
import type React from "react";
import { useDrainActions } from "@/components/DrainActions";
import {
  countMatchingWorkers,
  workerMatchesPattern,
} from "@/components/DrainsTable/utils";
import { usePlatformQueues } from "@/components/PlatformQueuesTable/bulkDrain";
import { findDrainTarget } from "@/components/PlatformQueuesTable/utils";
import WorkersDrainPanel from "@/components/WorkersDrainPanel";
import { useBbPortalMessage } from "@/context/MessageContext";
import type {
  PaginationInfo,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { Route, type WorkerListStatus } from "@/routes/scheduler.worker";
import themeStyles from "@/theme/theme.module.css";
import {
  type DrainPatternAction,
  dispatchDrainPattern,
  getDrainPattern,
  replacedValue,
  useDrainPattern,
} from "@/utils/drainPattern";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";
import { useOperatorMode } from "@/utils/operatorMode";
import WorkersTablePageSelector from "../WorkersTablePageSelector";
import WorkersTableTypeSelector from "../WorkersTableTypeSelector";
import getColumns from "./Columns";
import styles from "./index.module.css";

type Props = {
  workerStatusFilter: WorkerListStatus;
  sizeClassQueueName: SizeClassQueueNameSearch;
  data: WorkerState[] | undefined;
  paginationInfo: PaginationInfo | undefined;
  isLoading: boolean;
  pageSize: number;
};

const WorkersTable: React.FC<Props> = ({
  workerStatusFilter,
  sizeClassQueueName,
  data,
  paginationInfo,
  isLoading,
  pageSize,
}) => {
  const navigate = Route.useNavigate();
  const [operatorMode] = useOperatorMode();
  const { actions: drainActions, contextHolder } = useDrainActions(
    sizeClassQueueName,
    operatorMode,
  );
  const { data: platformQueues } = usePlatformQueues(
    drainActions !== undefined,
  );
  const queueWorkersCount =
    platformQueues &&
    findDrainTarget(platformQueues, sizeClassQueueName)?.workersCount;
  const { messageApi } = useBbPortalMessage();
  const { token } = antTheme.useToken();
  const workerIdPattern = useDrainPattern(sizeClassQueueName);
  const patternEmpty = Object.keys(workerIdPattern).length === 0;
  const match = data && {
    matching: countMatchingWorkers(
      data.map((worker) => worker.id),
      workerIdPattern,
    ),
    listed: data.length,
  };

  const updatePattern = (action: DrainPatternAction) => {
    const previous = replacedValue(getDrainPattern(sizeClassQueueName), action);
    dispatchDrainPattern(sizeClassQueueName, action);
    if (previous !== undefined && "value" in action) {
      messageApi.warning(
        `Replaced ${action.key}=${previous} with ${action.key}=${action.value}: a drain pattern holds one value per key.`,
      );
    }
  };

  const handleFilterChange = (value: WorkerListStatus) => {
    navigate({
      search: (prev) => ({
        ...prev,
        workerStatusFilter: value,
        cursor: undefined,
      }),
    });
  };

  const goToNextPage = () => {
    const lastId = data?.at(-1)?.id;
    if (!lastId) return;

    navigate({
      search: (prev) => ({
        ...prev,
        cursor: lastId,
      }),
    });
  };

  const goToFirstPage = () => {
    navigate({
      search: (prev) => ({
        ...prev,
        cursor: undefined,
      }),
    });
  };

  return (
    <>
      {contextHolder}
      <Space orientation="vertical" style={{ width: "100%" }}>
        <Row>
          <Flex style={{ width: "100%" }} justify="space-between" wrap>
            <WorkersTableTypeSelector
              workerStatusFilter={workerStatusFilter}
              setWorkerStatusFilter={handleFilterChange}
            />
            {paginationInfo && (
              <WorkersTablePageSelector
                paginationInfo={paginationInfo}
                goToFirstPage={goToFirstPage}
                goToNextPage={goToNextPage}
                pageSize={pageSize}
              />
            )}
          </Flex>
        </Row>
        {drainActions && (
          <Row>
            <WorkersDrainPanel
              workerIdPattern={workerIdPattern}
              match={match}
              onSet={(key, value) => updatePattern({ type: "set", key, value })}
              onRemove={(key) => updatePattern({ type: "remove", key })}
              onClear={() => updatePattern({ type: "clear" })}
              onDrainQueue={() => drainActions.drainQueue(queueWorkersCount)}
              onDrain={() =>
                drainActions.drain(workerIdPattern, {
                  match,
                  onDrained: () => updatePattern({ type: "clear" }),
                })
              }
            />
          </Row>
        )}
        <Row>
          <Table
            columns={getColumns(
              operatorMode,
              drainActions && {
                actions: drainActions,
                workerIdPattern,
                toggle: (key, value) =>
                  updatePattern({ type: "toggle", key, value }),
              },
            )}
            loading={isLoading}
            bordered={true}
            style={
              drainActions
                ? ({
                    width: "100%",
                    "--drain-match-bg": token.colorPrimaryBg,
                  } as React.CSSProperties)
                : { width: "100%" }
            }
            dataSource={data}
            size="small"
            rowClassName={(record) =>
              drainActions &&
              !patternEmpty &&
              workerMatchesPattern(record.id, workerIdPattern)
                ? `${themeStyles.compactTable} ${styles.drainMatch}`
                : themeStyles.compactTable
            }
            pagination={false}
            rowKey={(item) => Object.entries(item.id).sort().join(",")}
            locale={{
              emptyText:
                "No workers matching the given criteria found (that you have access to).",
            }}
          />
        </Row>
      </Space>
    </>
  );
};

export default WorkersTable;
