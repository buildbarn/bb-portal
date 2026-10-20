import { Button, Popconfirm, type TableColumnsType } from "antd";
import { DrainPattern } from "@/components/DrainActions";
import type { DrainState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { readableDurationFromDates } from "@/utils/time";
import { drainKey } from "./utils";

const getColumns = (
  canModify: boolean,
  onRemove: (workerIdPattern: DrainState["workerIdPattern"]) => void,
  removing: DrainState["workerIdPattern"] | undefined,
): TableColumnsType<DrainState> => [
  {
    key: "workerIdPattern",
    title: "Worker ID pattern",
    render: (_, record) => (
      <DrainPattern workerIdPattern={record.workerIdPattern} />
    ),
  },
  {
    key: "age",
    title: "Age",
    render: (_, record) =>
      record.createdTimestamp &&
      readableDurationFromDates(record.createdTimestamp, new Date(), {
        precision: 1,
        smallestUnit: "s",
      }),
  },
  ...(canModify
    ? [
        {
          key: "actions",
          title: "Actions",
          render: (_: unknown, record: DrainState) => (
            <Popconfirm
              title="Remove this drain?"
              description="Matching workers will pick up new operations again."
              okText="Remove"
              onConfirm={() => onRemove(record.workerIdPattern)}
            >
              <Button
                danger
                size="small"
                loading={
                  removing !== undefined &&
                  drainKey(removing) === drainKey(record.workerIdPattern)
                }
              >
                Remove
              </Button>
            </Popconfirm>
          ),
        },
      ]
    : []),
];

export default getColumns;
