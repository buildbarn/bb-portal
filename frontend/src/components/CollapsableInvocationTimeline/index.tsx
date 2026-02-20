import { ExclamationCircleOutlined } from "@ant-design/icons";
import { Collapse, Tooltip, Typography } from "antd";
import InvocationTimeline from "@/components/InvocationTimeline";
import type { GetBuildInvocationFragment } from "@/graphql/__generated__/graphql";

interface Props {
  invocations: GetBuildInvocationFragment[];
}

const CollapsableInvocationTimeline: React.FC<Props> = ({ invocations }) => {
  return (
    <Collapse
      bordered={false}
      items={[
        {
          key: 0,
          label: <Typography.Text strong>Invocation Timeline</Typography.Text>,
          extra: [
            <Tooltip
              key="0"
              title="Click on an entry to open the corresponding invocation page"
            >
              <ExclamationCircleOutlined />
            </Tooltip>,
          ],
          children: <InvocationTimeline invocations={invocations} />,
        },
      ]}
    />
  );
};

export default CollapsableInvocationTimeline;
