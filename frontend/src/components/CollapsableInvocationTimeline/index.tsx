import type { FindBuildFromUuidFragment } from "@/app/builds/[buildUUID]/[[...slugs]]/types";
import InvocationTimeline from "@/components/InvocationTimeline";
import { Collapse, Typography } from "antd";

interface Props {
  invocations: FindBuildFromUuidFragment[];
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
