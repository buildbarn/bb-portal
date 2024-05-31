import React from 'react';
import { Collapse, Descriptions, Typography } from 'antd';
import styles from './index.module.css';
import { ExitCode } from '@/graphql/__generated__/graphql';
import themeStyles from '@/theme/theme.module.css';

interface Props {
  invocationId: string;
  exitCode?: ExitCode | null;
}

const DebugInfo: React.FC<Props> = ({ invocationId, exitCode }) => {
  return (
    <Collapse
      bordered={false}
      className={themeStyles.collapse}
      items={[
        {
          key: 'debug',
          label: <Typography.Text strong>Bazel Information</Typography.Text>,
          children: (
            <Descriptions
              column={1}
              className={styles.descriptions}
              items={[
                {
                  key: 'invocation',
                  label: 'Invocation ID',
                  children: <Typography.Text code>{invocationId}</Typography.Text>,
                },
                {
                  key: 'code',
                  label: 'Exit Code',
                  children: exitCode ? (
                    <>
                      <Typography.Text code>{exitCode.code}</Typography.Text>{' '}
                      <Typography.Text code>{exitCode.name}</Typography.Text>
                    </>
                  ) : (
                    'None'
                  ),
                },
              ]}
            />
          ),
        },
      ]}
    />
  );
};

export default DebugInfo;
