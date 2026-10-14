import { Link } from "@tanstack/react-router";
import { Descriptions } from "antd";
import type React from "react";
import CodeText from "@/components/CodeText";
import PropertyTagList from "@/components/PropertyTagList";
import { InvocationChildrenFilter } from "@/routes/scheduler.invocation";
import {
  type InvocationNameSearch,
  invocationIdSearchToAny,
  invocationIdToString,
} from "@/utils/invocationName";

interface Props {
  invocationName: InvocationNameSearch;
}

const InvocationInfo: React.FC<Props> = ({ invocationName }) => {
  const { sizeClassQueueName, ids } = invocationName;
  const ancestorLink = (depth: number, children: React.ReactNode) => (
    <Link
      to="/scheduler/invocation"
      search={{
        invocationName: { sizeClassQueueName, ids: ids.slice(0, depth) },
        childrenFilter: InvocationChildrenFilter.ACTIVE,
        cursor: undefined,
      }}
    >
      {children}
    </Link>
  );

  return (
    <Descriptions column={1} bordered>
      <Descriptions.Item label="Instance name prefix">
        {sizeClassQueueName.platformQueueName.instanceNamePrefix}
      </Descriptions.Item>
      <Descriptions.Item label="Platform properties">
        <PropertyTagList
          propertyList={
            sizeClassQueueName.platformQueueName.platform.properties
          }
        />
      </Descriptions.Item>
      <Descriptions.Item label="Size class">
        {sizeClassQueueName.sizeClass}
      </Descriptions.Item>
      <Descriptions.Item label="Invocation IDs">
        <ul>
          <li>
            {ids.length
              ? ancestorLink(0, "Root invocation")
              : "Root invocation"}
          </li>
          {ids.map((id, i) => {
            const text = (
              <CodeText style={{ wordBreak: "break-all" }}>
                {invocationIdToString(invocationIdSearchToAny(id))}
              </CodeText>
            );
            return (
              // biome-ignore lint/suspicious/noArrayIndexKey: The same ID may appear at multiple depths
              <li key={`${i}:${id.typeUrl}:${id.value}`}>
                {i < ids.length - 1 ? ancestorLink(i + 1, text) : text}
              </li>
            );
          })}
        </ul>
      </Descriptions.Item>
    </Descriptions>
  );
};

export default InvocationInfo;
