/* eslint-disable react/no-array-index-key */
import React from 'react';
import { Button, Collapse, CollapseProps } from 'antd';
import CopyTextButton from '@/components/CopyTextButton';
import PortalAlert from '@/components/PortalAlert';
import themeStyles from '@/theme/theme.module.css';
import { ProblemInfoFragment } from "@/graphql/__generated__/graphql";
import BuildProblem, { BuildProblemLabel } from "@/components/Problems/BuildProblem";
import { ExclamationCircleFilled } from '@ant-design/icons';

interface Props {
  problems?: ProblemInfoFragment[];
}

export const CopyAllProblemLabels: React.FC<{ problems: ProblemInfoFragment[] }> = ({ problems }) => {
  // NOTE: Simplified since ProgressProblem has an '' label.
  return <CopyTextButton buttonText="Copy Problems" copyText={problems.map(problem => problem.label).join(' ')} />;
};

const BuildProblems: React.FC<Props> = ({ problems }) => {
  if (!problems || problems?.length === 0) {
    return (
      <div>
        <PortalAlert
          message="There are reported errors for this invocation."
          type="warning"
          showIcon
        />

        Click here to view error details
      </div>
    );
  }

  const progressID = problems.find(problem => problem.__typename === 'ProgressProblem');
  const defaultActiveKey = progressID?.id ?? (problems.length === 1 ? problems[0].id : undefined);

  // Map all problems to build problem components.
  const items: CollapseProps['items'] = problems.map(problem => {
    return {
      key: problem.id,
      label: <BuildProblemLabel problem={problem} />,
      children: <BuildProblem problem={problem} />,
    };
  });

  // If there is only one problem, expand it.
  return (
    <Collapse
      items={items}
      defaultActiveKey={defaultActiveKey}
      bordered={false}
      destroyInactivePanel
      className={themeStyles.collapse}
    />
  );
};

export default BuildProblems;
