import { MoreOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { Button, Dropdown, Space } from "antd";
import type React from "react";
import { useDrainActions } from "@/components/DrainActions";
import { WorkerListStatus } from "@/routes/scheduler.worker";
import type { SizeClassQueueNameSearch } from "@/utils/invocationName";

interface Props {
  sizeClassQueueName: SizeClassQueueNameSearch;
  drainsCount: number;
  workersCount: number;
}

const DrainsCell: React.FC<Props> = ({
  sizeClassQueueName,
  drainsCount,
  workersCount,
}) => {
  const { actions, contextHolder } = useDrainActions(sizeClassQueueName, true);

  return (
    <>
      {contextHolder}
      <Space size="small">
        <Link to="/scheduler/drain" search={{ sizeClassQueueName }}>
          {drainsCount}
        </Link>
        {actions && (
          <Dropdown
            trigger={["click"]}
            menu={{
              items: [
                {
                  key: "drainQueue",
                  danger: true,
                  label: "Drain all workers in this size class queue",
                  onClick: () => actions.drainQueue(workersCount),
                },
                {
                  key: "workers",
                  label: (
                    <Link
                      to="/scheduler/worker"
                      search={{
                        workerStatusFilter: WorkerListStatus.ALL,
                        sizeClassQueueName,
                        cursor: undefined,
                      }}
                    >
                      Build a drain on the workers page
                    </Link>
                  ),
                },
              ],
            }}
          >
            <Button
              size="small"
              type="text"
              icon={<MoreOutlined />}
              aria-label="Drain actions"
            />
          </Dropdown>
        )}
      </Space>
    </>
  );
};

export default DrainsCell;
