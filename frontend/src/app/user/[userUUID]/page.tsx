"use client";

import { CalendarFilled } from "@ant-design/icons";
import { useQuery } from "@apollo/client";
import { Flex, Spin } from "antd";
import type React from "react";
import Content from "@/components/Content";
import PortalAlert from "@/components/PortalAlert";
import PortalCard from "@/components/PortalCard";
import UserView from "@/components/UserView";
import { getFragmentData } from "@/graphql/__generated__";
import {
  BazelInvocationOrderField,
  OrderDirection,
} from "@/graphql/__generated__/graphql";
import GET_AUTHENTICATED_USER_BY_UUID, {
  AUTHENTICATED_USER_NODE_FRAGMENT,
} from "./index.graphql";
import styles from "./index.module.css";

interface PageParams {
  params: {
    userUUID: string;
  };
}

const Page: React.FC<PageParams> = ({ params }) => {
  const { loading, data, error } = useQuery(GET_AUTHENTICATED_USER_BY_UUID, {
    variables: {
      userUUID: params.userUUID,
      bazelInvocationsOrderBy: {
        field: BazelInvocationOrderField.StartedAt,
        direction: OrderDirection.Desc,
      },
    },
    fetchPolicy: "network-only",
  });

  if (loading)
    return (
      <Flex justify="center">
        <Spin />
      </Flex>
    );

  if (error)
    return (
      <PortalAlert
        type="error"
        message={`There was a problem communicating with the backend server: ${error?.message}`}
        showIcon
        className={styles.alert}
      />
    );

  const authenticatedUser = getFragmentData(
    AUTHENTICATED_USER_NODE_FRAGMENT,
    data?.getAuthenticatedUser,
  );

  return (
    <Content
      content={
        <PortalCard
          icon={<CalendarFilled />}
          titleBits={[
            <span key="title">
              User {authenticatedUser?.displayName || params.userUUID}
            </span>,
          ]}
        >
          <UserView authenticatedUser={authenticatedUser} />
        </PortalCard>
      }
    />
  );
};

export default Page;
