import {Descriptions, Empty, Space, Typography} from "antd";
import styles from './index.module.css';
import React from "react";
import {DescriptionsItemType} from "antd/es/descriptions";
import {ProgressProblem, TargetProblem, TestProblem} from "@/graphql/__generated__/graphql";
import {LogViewerCard} from "@/components/LogViewer";
import themeStyles from '@/theme/theme.module.css';
import {Problem} from "@/components/Problems/types";
import {TestStatus} from "@/components/Problems/TestStatus";
import TestResultContainer from "@/components/Problems/TestResultContainer";
import ActionProblemContainer from "@/components/Problems/ActionProblemContainer";


const BuildLabel: React.FC<{ label: string; copyable?: boolean; copyText?: string }> = ({
  label,
  copyable,
  copyText,
}) => {
  return (
    <div className={styles.labelContainer}>
      <span className={styles.labelText}>{label}</span>
      {copyable ? (
        <span className={styles.copyIcon}>
          <Typography.Text copyable={{ text: copyText ?? label }}></Typography.Text>
        </span>
      ) : null}
    </div>
  );
};

export const BuildProblemLabel: React.FC<Props> = ({ problem }) => {
  // eslint-disable-next-line no-underscore-dangle
  switch (problem.__typename) {
    case 'ActionProblem':
      return <BuildLabel label={`${problem.type} @ ${problem.label}`} copyable copyText={problem.label} />;
    case 'ProgressProblem':
      return <BuildLabel label="Console log with errors" />;
    case 'TestProblem':
    case 'TargetProblem':
      return <BuildLabel label={problem.label} copyable />;
    default:
      return null;
  }
};

export const SectionWithTestStatus: React.FC<{
  testProblem: TestProblem;
  extraItems?: DescriptionsItemType[];
  children?: React.ReactNode;
}> = ({ testProblem, extraItems, children }) => {
  return (
    <>
      <Descriptions
        column={1}
        items={[
          {
            key: 'status',
            label: 'Status',
            children: <TestStatus status={testProblem.status} />,
          },
          ...(extraItems ?? []),
        ]}
      />
      {children}
    </>
  );
};

interface TestProblemPanelProps {
  testProblem: TestProblem;
}

const TestProblemPanel: React.FC<TestProblemPanelProps> = ({ testProblem }) => {
  if (testProblem.results.length > 0) {
    return (
      <TestResultContainer id={testProblem.results[0].id} problemLabel={testProblem.label} testProblem={testProblem} />
    );
  }

  return (
    <SectionWithTestStatus testProblem={testProblem}>
      <Empty description="No test result available" />
    </SectionWithTestStatus>
  );
};

interface TargetProblemPanelProps {
  targetProblem: TargetProblem;
}

const TargetProblemPanel: React.FC<TargetProblemPanelProps> = ({ targetProblem }) => {
  return (
    <Space direction="vertical" size="middle" className={themeStyles.space}>
      <LogViewerCard log={null} />
    </Space>
  );
};

interface ProgressProblemPanelProps {
  progressProblem: ProgressProblem;
}

const ProgressProblemPanel: React.FC<ProgressProblemPanelProps> = ({ progressProblem }) => {
  return <LogViewerCard log={progressProblem.output} />;
};

interface Props {
  problem: Problem;
}

const BuildProblem: React.FC<Props> = ({ problem }) => {
  // eslint-disable-next-line no-underscore-dangle
  switch (problem.__typename) {
    case 'ActionProblem':
      return <ActionProblemContainer id={problem.id} />;
    case 'TestProblem':
      return <TestProblemPanel testProblem={problem as TestProblem} />;
    case 'TargetProblem':
      return <TargetProblemPanel targetProblem={problem} />;
    case 'ProgressProblem':
      return <ProgressProblemPanel progressProblem={problem} />;
    default:
      return null;
  }
};

export default BuildProblem;
