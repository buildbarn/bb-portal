import React, { RefAttributes } from 'react';
import { Card, CardProps } from 'antd';
import { AnsiUp } from 'ansi_up';
import linkifyHtml from 'linkify-html';
import { JSX } from 'react/jsx-runtime';
import styles from './index.module.css';
import PortalAlert from '@/components/PortalAlert';
import IntrinsicAttributes = JSX.IntrinsicAttributes;

const ansi = new AnsiUp();
//large logs were causing a stack overflow, so setting a configurable max length and tailing the logs for that length
const MAX_LOG_LENGTH = 50000

const FILE_EXTENSIONS_IGNORE = ['.py', '.so'];

interface Props {
  log?: string | null;
  copyable?: boolean;
}
const LogViewer: React.FC<Props> = ({ log }) => {
  if (!log) {
    return (
      <PortalAlert message="There is no information to display" type="warning" showIcon className={styles.alert} />
    );
  }
  if (log.length > MAX_LOG_LENGTH) {
    log = log.slice(0, MAX_LOG_LENGTH / 2) + "\n\n...\n\n**************LOG CONTENT TRUNCATED********************\n\n...\n\n" + log.slice(log.length - (MAX_LOG_LENGTH / 2), log.length)
  }
  const innerHTML =
    linkifyHtml(ansi.ansi_to_html(log), {
      target: '_blank',
      validate: {
        url: url => {
          try {
            new URL(url);
            return !FILE_EXTENSIONS_IGNORE.some(fileExtension => {
              const ignoredFileExtensionRegExp = new RegExp(`\\w\\${fileExtension}.*`);
              return ignoredFileExtensionRegExp.test(url);
            });
          } catch (TypeError) {
            return false;
          }
        },
      },
    })
  return <pre dangerouslySetInnerHTML={{ __html: innerHTML }} />;
};

type LogViewerCardProps = Props & IntrinsicAttributes & CardProps & RefAttributes<HTMLDivElement>;

export const LogViewerCard: React.FC<LogViewerCardProps> = ({ log, bordered, copyable, ...props }) => {
  return (
    <Card bordered={false} {...props}>
      <LogViewer log={log} copyable={copyable} />
    </Card>
  );
};

export default LogViewer;
