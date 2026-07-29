import { BuildOutlined } from "@ant-design/icons";
import { Link, Outlet } from "@tanstack/react-router";
import { Typography } from "antd";
import { useMemo } from "react";
import styles from "@/components/AppBar/index.module.css";
import { BazelInvocationTabBar } from "@/components/BazelInvocationTabBar";
import { InvocationResultTag } from "@/components/InvocationResultTag";
import { PortalCard } from "@/components/PortalCard";
import PortalDuration from "@/components/PortalDuration";
import ProfileDropdown from "@/components/ProfileDropdown";
import UserStatusIndicator from "@/components/UserStatusIndicator";
import { getFragmentData } from "@/graphql/__generated__";
import type { BazelInvocationCommonFragment } from "@/graphql/__generated__/graphql";
import { FILE_DETAILS_FRAGMENT } from "@/types/GraphqlFileFragment";

const getTitleBits = (
  invocation: BazelInvocationCommonFragment,
): React.ReactNode[] => {
  const { invocationID, authenticatedUser, username } = invocation;

  const titleBits: React.ReactNode[] = [];
  if (username && username !== "")
    titleBits.push(
      <span key="username">
        User:{" "}
        <Typography.Text type="secondary" className={styles.normalWeight}>
          <UserStatusIndicator
            authenticatedUser={authenticatedUser}
            username={username}
            showIcon
          />
        </Typography.Text>
      </span>,
    );
  if (invocationID && invocationID !== "")
    titleBits.push(
      <span key="invocationID">
        Invocation ID:{" "}
        <Typography.Text
          type="secondary"
          className={styles.normalWeight}
          copyable={{ text: invocationID ?? "Copy" }}
        >
          {invocationID}
        </Typography.Text>{" "}
      </span>,
    );
  titleBits.push(
    <InvocationResultTag
      key="result"
      exitCodeName={invocation.exitCodeName || undefined}
      timeSinceLastConnectionMillis={
        invocation.connectionMetadata?.timeSinceLastConnectionMillis ||
        undefined
      }
    />,
  );
  return titleBits;
};

const getExtraBits = (
  invocation: BazelInvocationCommonFragment,
): React.ReactNode[] => {
  const { invocationID, build, profile } = invocation;

  const extraBits: React.ReactNode[] = [];
  if (build?.buildUUID) {
    extraBits.push(
      <span key="build">
        Build{" "}
        <Link to={`/builds/$buildUUID`} params={{ buildUUID: build.buildUUID }}>
          {build.buildUUID}
        </Link>
      </span>,
    );
  }
  extraBits.push(
    <PortalDuration
      key="duration"
      from={invocation.startedAt || undefined}
      to={
        invocation.endedAt
          ? invocation.endedAt
          : invocation.connectionMetadata?.connectionLastOpenAt
      }
      includeIcon
      formatConfig={{ smallestUnit: "s" }}
    />,
  );
  if (profile) {
    const parsedProfile = getFragmentData(FILE_DETAILS_FRAGMENT, profile);
    extraBits.push(
      <ProfileDropdown profile={parsedProfile} invocationID={invocationID} />,
    );
  }
  return extraBits;
};

interface Props {
  invocation: BazelInvocationCommonFragment;
}

export const BazelInvocationDetailsPage: React.FC<Props> = ({ invocation }) => {
  const titleBits = useMemo(() => getTitleBits(invocation), [invocation]);
  const extraBits = useMemo(() => getExtraBits(invocation), [invocation]);

  return (
    <PortalCard
      icon={<BuildOutlined />}
      titleBits={titleBits}
      extraBits={extraBits}
    >
      <BazelInvocationTabBar invocation={invocation} />
      <Outlet />
    </PortalCard>
  );
};
