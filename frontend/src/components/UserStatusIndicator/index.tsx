import { UserOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import { theme } from "antd";
import type { BazelInvocationNodeFragment } from "@/graphql/__generated__/graphql";

const { useToken } = theme;
interface Props {
  authenticatedUser?: BazelInvocationNodeFragment["authenticatedUser"];
  showIcon?: boolean;
  username?: string;
}

const UserStatusIndicator: React.FC<Props> = ({
  authenticatedUser,
  showIcon,
  username,
}) => {
  const { token } = useToken();
  if (authenticatedUser) {
    return (
      <>
        {showIcon && (
          <UserOutlined
            style={{ color: token.colorText, marginRight: "5px" }}
            title="This user is authenticated"
          />
        )}
        <span style={{ color: token.colorLink }}>
          <Link
            to="/users/$userUUID"
            params={{ userUUID: authenticatedUser.userUUID }}
          >
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
      {showIcon && (
        <UserOutlined
          style={{ color: token.red, marginRight: "5px" }}
          title="This user is unauthenticated"
        />
      )}
      {username || <i>No display name</i>}
    </>
  );
};

export default UserStatusIndicator;
