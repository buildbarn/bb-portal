import React, { RefAttributes, useState } from "react";
import { Card, CardProps } from "antd";
import { AnsiUp } from "ansi_up";
import linkifyHtml from "linkify-html";
import { JSX } from "react/jsx-runtime";
import styles from "./index.module.css";
import PortalAlert from "@/components/PortalAlert";
import IntrinsicAttributes = JSX.IntrinsicAttributes;
import { GET_BUILD_LOGS } from "./graphql";
import { GetBuildLogsQueryVariables } from "@/graphql/__generated__/graphql";
import { useQuery } from "@apollo/client";

const ansi = new AnsiUp();
const MAX_LOG_LENGTH = 50000;

const FILE_EXTENSIONS_IGNORE = [".py", ".so"];

interface Props {
  invocationId?: string | null;
  log?: string | null;
  copyable?: boolean;
}

const LogViewer: React.FC<Props> = ({ log, invocationId }) => {
  if (!log) {
    return (
      <PortalAlert
        message="There is no log information to display"
        type="warning"
        showIcon
        className={styles.alert}
      />
    );
  } else {
    const innerHTML = ansi.ansi_to_html(log);
    return <pre dangerouslySetInnerHTML={{ __html: innerHTML }} />;
  }
};

type LogViewerCardProps = Props &
  IntrinsicAttributes &
  CardProps &
  RefAttributes<HTMLDivElement>;

export const LogViewerCard: React.FC<LogViewerCardProps> = ({
  log,
  invocationId,
  bordered,
  copyable,
  ...props
}) => {
  return (
    <Card bordered={false} {...props}>
      <LogViewer log={log} invocationId={invocationId} copyable={copyable} />
    </Card>
  );
};

export default LogViewer;
