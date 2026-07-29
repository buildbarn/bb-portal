import { Link } from "@tanstack/react-router";
import { Space, Typography } from "antd";
import type React from "react";
import themeStyles from "@/theme/theme.module.css";

interface Props {
  type?: string;
  details?: string;
  showUnauthenticatedMessage?: boolean;
}

export const NotFoundPage: React.FC<Props> = ({
  type,
  details,
  showUnauthenticatedMessage = false,
}) => {
  return (
    <Space
      direction="vertical"
      size="large"
      className={themeStyles.errorPageContainer}
    >
      <Typography.Title>
        The {type || "page"} you’re looking for can’t be found.
      </Typography.Title>
      {showUnauthenticatedMessage && (
        <Typography.Text>
          Either the {type} doesn't exist or you don't have access to it.
        </Typography.Text>
      )}
      {details && <Typography.Text>{details}</Typography.Text>}
      <Link to="/">Go Back Home</Link>
    </Space>
  );
};
