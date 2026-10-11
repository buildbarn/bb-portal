import { useQueryClient } from "@tanstack/react-query";
import { Modal, Space, Table, Tag, Typography } from "antd";
import type React from "react";
import {
  invalidateDrainQueries,
  useAddDrain,
  useCanModifyDrains,
} from "@/components/DrainsTable/hooks";
import {
  patternToPropertyList,
  patternToString,
} from "@/components/DrainsTable/utils";
import type { DrainTarget } from "@/components/PlatformQueuesTable/utils";
import PropertyTagList from "@/components/PropertyTagList";
import WorkersInfo from "@/components/WorkersInfo";
import { useBbPortalMessage } from "@/context/MessageContext";
import { buildQueueStateClient } from "@/grpc/buildQueueStateClient";
import { grpcErrorText } from "@/utils/grpcStatus";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";
import styles from "./index.module.css";

export interface ListedMatch {
  matching: number;
  listed: number;
}

export interface DrainActions {
  drain: (
    workerIdPattern: Record<string, string>,
    options?: { match?: ListedMatch; onDrained?: () => void },
  ) => void;
  drainQueue: (workersCount: number | undefined) => void;
}

export const DrainPattern: React.FC<{
  workerIdPattern: Record<string, string>;
}> = ({ workerIdPattern }) => {
  const propertyList = patternToPropertyList(workerIdPattern);
  return propertyList.length ? (
    <PropertyTagList propertyList={propertyList} />
  ) : (
    <Tag>all workers</Tag>
  );
};

export const useDrainActions = (
  sizeClassQueueName: SizeClassQueueNameSearch,
  enabled: boolean,
): { actions: DrainActions | undefined; contextHolder: React.ReactNode } => {
  const [modal, contextHolder] = Modal.useModal();
  const { messageApi } = useBbPortalMessage();
  const canModify = useCanModifyDrains(
    sizeClassQueueName.platformQueueName.instanceNamePrefix,
    enabled,
  );
  const addDrain = useAddDrain();

  const submit = (workerIdPattern: Record<string, string>) =>
    addDrain.mutateAsync({ sizeClassQueueName, workerIdPattern }).then(
      () => {
        messageApi.success(
          Object.keys(workerIdPattern).length
            ? `Drained workers matching ${patternToString(workerIdPattern)}.`
            : "Drained all workers in the size class queue.",
        );
        return true;
      },
      (error: Error) => {
        messageApi.error(`Failed to add drain: ${grpcErrorText(error)}`);
        return false;
      },
    );

  const queueInfo = (
    <WorkersInfo sizeClassQueueName={sizeClassQueueName} size="small" />
  );

  const drain: DrainActions["drain"] = (workerIdPattern, options) =>
    modal.confirm({
      title: "Drain matching workers?",
      width: 560,
      okText: "Drain",
      okButtonProps: { danger: true },
      content: (
        <Space
          orientation="vertical"
          className={styles.content}
          style={{ width: "100%" }}
        >
          <Typography.Text>
            Every worker in this size class queue whose ID contains all of
          </Typography.Text>
          <DrainPattern workerIdPattern={workerIdPattern} />
          <Typography.Text>
            finishes its current operation and receives no new ones until the
            drain is removed.
          </Typography.Text>
          {queueInfo}
          {options?.match && (
            <Typography.Text type="secondary">
              {options.match.matching} of the {options.match.listed} workers on
              this page match.
            </Typography.Text>
          )}
        </Space>
      ),
      onOk: () =>
        submit(workerIdPattern).then((ok) => ok && options?.onDrained?.()),
    });

  const drainQueue: DrainActions["drainQueue"] = (workersCount) =>
    modal.confirm({
      title: "Drain every worker in this size class queue?",
      width: 560,
      okText:
        workersCount === undefined
          ? "Drain all workers"
          : `Drain all ${workersCount} workers`,
      okButtonProps: { danger: true },
      content: (
        <Space
          orientation="vertical"
          className={styles.content}
          style={{ width: "100%" }}
        >
          <Typography.Text type="danger" strong>
            This drains{" "}
            {workersCount === undefined
              ? "every worker"
              : `all ${workersCount} workers`}{" "}
            in the size class queue below. They finish their current operation
            and then receive no new ones, so queued and newly submitted
            operations for this queue do not run until the drain is removed.
          </Typography.Text>
          {queueInfo}
        </Space>
      ),
      onOk: () => submit({}),
    });

  return {
    actions: enabled && canModify === true ? { drain, drainQueue } : undefined,
    contextHolder,
  };
};

export const drainTargetLabel = ({ sizeClassQueueName }: DrainTarget) => {
  const { instanceNamePrefix, platform } = sizeClassQueueName.platformQueueName;
  return `"${instanceNamePrefix}" [${platform.properties
    .map(({ name, value }) => `${name}=${value}`)
    .join(", ")}] size class ${sizeClassQueueName.sizeClass}`;
};

export const useBulkDrain = (
  isAllowed: (instanceNamePrefix: string) => boolean | undefined,
): {
  drainTargets: (targets: DrainTarget[], description: string) => void;
  contextHolder: React.ReactNode;
} => {
  const [modal, contextHolder] = Modal.useModal();
  const { messageApi } = useBbPortalMessage();
  const queryClient = useQueryClient();

  const run = async (targets: DrainTarget[], description: string) => {
    const failures: { target: DrainTarget; error: Error }[] = [];
    for (const target of targets) {
      try {
        await buildQueueStateClient.addDrain({
          sizeClassQueueName: target.sizeClassQueueName,
          workerIdPattern: {},
        });
      } catch (error) {
        failures.push({ target, error: error as Error });
      }
    }
    invalidateDrainQueries(queryClient);
    const summary = `Drained ${targets.length - failures.length} of ${targets.length} queues ${description}`;
    if (failures.length === 0) {
      messageApi.success(`${summary}.`);
      return;
    }
    messageApi.error(
      `${summary}; ${failures
        .map(
          ({ target, error }) =>
            `size class ${target.sizeClassQueueName.sizeClass} failed: ${grpcErrorText(error)}`,
        )
        .join("; ")}`,
      10,
    );
    modal.error({
      title: summary,
      width: 720,
      content: (
        <ul>
          {failures.map(({ target, error }) => (
            <li key={drainTargetLabel(target)}>
              {drainTargetLabel(target)}: {grpcErrorText(error)}
            </li>
          ))}
        </ul>
      ),
    });
  };

  const drainTargets = (targets: DrainTarget[], description: string) => {
    const allowed = targets.filter(({ sizeClassQueueName }) =>
      isAllowed(sizeClassQueueName.platformQueueName.instanceNamePrefix),
    );
    const skipped = targets.length - allowed.length;
    const workers = allowed.reduce((sum, t) => sum + t.workersCount, 0);
    modal.confirm({
      title: `Drain all workers in every queue ${description}?`,
      width: 880,
      okText: `Drain ${allowed.length} queues (${workers} workers)`,
      okButtonProps: { danger: true, disabled: allowed.length === 0 },
      content: (
        <Space
          orientation="vertical"
          className={styles.content}
          style={{ width: "100%" }}
        >
          <Typography.Text type="danger" strong>
            This adds a drain matching every worker to each of the{" "}
            {allowed.length} size class queues below, affecting {workers}{" "}
            workers. Running operations finish, but queued and newly submitted
            operations for these queues do not run until the drains are removed.
          </Typography.Text>
          {skipped > 0 && (
            <Typography.Text type="warning">
              {skipped} of the {targets.length} queues are skipped because you
              are not allowed to modify their drains.
            </Typography.Text>
          )}
          <Table
            size="small"
            bordered
            pagination={false}
            dataSource={targets}
            rowKey={drainTargetLabel}
            columns={[
              {
                key: "instanceNamePrefix",
                title: "Instance name prefix",
                render: (_, t) =>
                  t.sizeClassQueueName.platformQueueName.instanceNamePrefix,
              },
              {
                key: "platform",
                title: "Platform properties",
                render: (_, t) => (
                  <PropertyTagList
                    propertyList={
                      t.sizeClassQueueName.platformQueueName.platform.properties
                    }
                  />
                ),
              },
              {
                key: "sizeClass",
                title: "Size class",
                render: (_, t) => t.sizeClassQueueName.sizeClass,
              },
              {
                key: "workers",
                title: "Workers",
                render: (_, t) => t.workersCount,
              },
              {
                key: "drains",
                title: "Existing drains",
                render: (_, t) => t.drainsCount,
              },
              {
                key: "action",
                title: "Action",
                className: styles.nowrap,
                render: (_, t) =>
                  isAllowed(
                    t.sizeClassQueueName.platformQueueName.instanceNamePrefix,
                  ) ? (
                    <Tag color="red">drain</Tag>
                  ) : (
                    <Tag>skipped: not allowed</Tag>
                  ),
              },
            ]}
          />
        </Space>
      ),
      onOk: () => run(allowed, description),
    });
  };

  return { drainTargets, contextHolder };
};
