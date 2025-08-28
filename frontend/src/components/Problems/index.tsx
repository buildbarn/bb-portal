/* eslint-disable react/no-array-index-key */
import React, { useState } from 'react';
import { Button, Collapse, CollapseProps } from 'antd';
import CopyTextButton from '@/components/CopyTextButton';
import PortalAlert from '@/components/PortalAlert';
import themeStyles from '@/theme/theme.module.css';
import { ProblemInfoFragment } from "@/graphql/__generated__/graphql";
import BuildProblem, { BuildProblemLabel } from "@/components/Problems/BuildProblem";
import { ExclamationCircleFilled } from '@ant-design/icons';
import { useQuery } from '@apollo/client';
import { GET_PROBLEM_DETAILS, PROBLEM_DETAILS_FRAGMENT, PROBLEM_INFO_FRAGMENT } from '@/app/bazel-invocations/[invocationID]/index.graphql';
import { getFragmentData } from '@/graphql/__generated__';
import { domainToASCII } from 'url';

interface Props {
  invocationId: string;
  instanceName: string | undefined;
  //problems?: ProblemInfoFragment[];
}


export const CopyAllProblemLabels: React.FC<{ problems: ProblemInfoFragment[] }> = ({ problems }) => {
  // NOTE: Simplified since ProgressProblem has an '' label.
  return <CopyTextButton buttonText="Copy Problems" copyText={problems.map(problem => problem.label).join(' ')} />;
};

const BuildProblems: React.FC<Props> = ({ invocationId, instanceName }) => {

  var { loading, data, previousData, error } = useQuery(GET_PROBLEM_DETAILS, {
    variables: {
      invocationID: invocationId
    }, fetchPolicy: 'cache-and-network'
  });

  var activeData = loading ? previousData : data;
  var problems: ProblemInfoFragment[] | undefined = []

  if (error) {
    problems = []
  } else {
    const invocation = getFragmentData(PROBLEM_DETAILS_FRAGMENT, activeData?.bazelInvocation)
    problems = invocation?.problems.map(p => getFragmentData(PROBLEM_INFO_FRAGMENT, p))
    if (!problems || problems?.length === 0) {
      return (
        <PortalAlert
          message="There is no debug information for this invocation."
          type="warning"
          showIcon
        />
      );
    }
  }

  const progressID = problems.find(problem => problem.__typename === 'ProgressProblem');
  const defaultActiveKey = progressID?.id ?? (problems.length === 1 ? problems[0].id : undefined);

  // Map all problems to build problem components.
  const items: CollapseProps['items'] = problems.map(problem => {
    return {
      key: problem.id,
      label: <BuildProblemLabel problem={problem} instanceName={instanceName}/>,
      children: <BuildProblem problem={problem} instanceName={instanceName} />,
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
