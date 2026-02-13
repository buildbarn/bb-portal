import { Collapse, Typography } from "antd";
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
          children: <InvocationTimeline invocations={invocations} />,
        },
      ]}
    />
  );
};

export default CollapsableInvocationTimeline;
