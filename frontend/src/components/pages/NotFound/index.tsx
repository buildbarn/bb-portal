import { Link } from "@tanstack/react-router";
import { Space, Typography } from "antd";
import type React from "react";
import Content from "@/components/Content";
import themeStyles from "@/theme/theme.module.css";

interface Props {
  type?: string;
  showUnauthenticatedMessage?: boolean;
}

export const NotFoundPage: React.FC<Props> = ({
  type,
  showUnauthenticatedMessage = false,
}) => {
  return (
    <Content
      content={
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
          <Link to="/">Go Back Home</Link>
        </Space>
      }
    />
  );
};
