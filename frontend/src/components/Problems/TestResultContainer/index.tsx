import { TestProblem } from "@/graphql/__generated__/graphql";
import { Space } from "antd";
import React from "react";
import ErrorAlert from "@/components/ErrorAlert";
import themeStyles from '@/theme/theme.module.css';
import { SectionWithTestStatus } from "@/components/Problems/BuildProblem";
import LogOutput from "@/components/Problems/LogOutput";

interface Props {
  id: string;
  problemLabel: string;
  testProblem: TestProblem;
  instanceName: string | undefined;
}


const TestResultContainer: React.FC<Props> = ({ id, problemLabel, testProblem, instanceName}) => {
  const testResult = testProblem.results.find(r => r.id == id)
  if (!testResult) {
    return (
      <>
        <SectionWithTestStatus testProblem={testProblem} />
        <ErrorAlert error={new Error('Expected test result but server returned something else')} />
      </>
    );
  }

  var contents = <LogOutput blobReference={testResult.actionLogOutput} instanceName={instanceName}/>


  return (
    <>
      <SectionWithTestStatus testProblem={testProblem} />
      {/* Display spin behind the actions, making UI stable when query is being executed. */}
      <Space direction="vertical" size="middle" className={themeStyles.space}>
        {contents}
      </Space>
    </>
  );
};

export default TestResultContainer;
