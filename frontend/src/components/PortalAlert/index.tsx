import { Alert, Typography } from "antd";
import type { AlertProps } from "antd/lib/alert";
import type React from "react";
import styles from "./index.module.css";

type Props = AlertProps;

const PortalAlert: React.FC<Props> = ({ className, message, ...props }) => {
  return (
    <Alert
      className={[styles.alert, className].join(" ")}
      message={
        message && props.description ? (
          <Typography.Title level={5}>{message}</Typography.Title>
        ) : (
          message
        )
      }
      {...props}
    />
  );
};

export default PortalAlert;
