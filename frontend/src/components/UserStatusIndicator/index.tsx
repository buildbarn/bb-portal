import { UserOutlined } from "@ant-design/icons";
import { Tooltip, theme } from "antd";
import Link from "next/link";
import type { BazelInvocationNodeFragment } from "@/graphql/__generated__/graphql";

const { useToken } = theme;
interface Props {
  authenticatedUser: BazelInvocationNodeFragment["authenticatedUser"];
  user: BazelInvocationNodeFragment["user"];
}

const UserStatusIndicator: React.FC<Props> = ({ authenticatedUser, user }) => {
  const { token } = useToken();
  if (authenticatedUser) {
    return (
      <>
        <Tooltip title="This user is authenticated">
          <UserOutlined style={{ color: token.colorText }} />
        </Tooltip>{" "}
        <span style={{ color: token.colorLink }}>
          <Link href={`/user/${authenticatedUser.userUUID}`}>
            {authenticatedUser?.displayName !== "" ? (
              authenticatedUser?.displayName
            ) : (
              <i>No display name</i>
            )}
          </Link>
        </span>
      </>
    );
  }

  return (
    <>
      <Tooltip title="This user is unauthenticated">
        <UserOutlined style={{ color: "red" }} />
      </Tooltip>{" "}
      {user?.LDAP}
    </>
  );
};

export default UserStatusIndicator;
