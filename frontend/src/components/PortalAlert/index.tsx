import React from 'react';
import { Alert, Typography } from 'antd';
import { AlertProps } from 'antd/lib/alert';
import styles from './index.module.css';
import IntrinsicAttributes = JSX.IntrinsicAttributes;

type Props = IntrinsicAttributes & AlertProps;

const PortalAlert: React.FC<Props> = ({ className, message, ...props }) => {
    return (
        <Alert
            className={[styles.alert, className].join(' ')}
            message={message && props.description ? <Typography.Title level={5}>{message}</Typography.Title> : message}
            {...props}
        />
    );
};

export default PortalAlert;
