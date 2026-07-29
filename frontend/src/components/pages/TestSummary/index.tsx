import { ExperimentFilled } from "@ant-design/icons";
import { Descriptions, Space, Statistic, Table, Typography } from "antd";
import type React from "react";
import { useMemo } from "react";
import {
  type CacheLocation,
  CacheLocationTag,
  cacheLocationFromBooleans,
  cacheLocationFromTestResults,
} from "@/components/CacheLocationTag";
import { LinkButton } from "@/components/LinkButton";
import { CasGqlFileViewer } from "@/components/LogViewer/casGqlFileViewer";
import PortalCard from "@/components/PortalCard";
import { TestStatusTag } from "@/components/TestStatusTag";
import { getFragmentData } from "@/graphql/__generated__";
import type {
  FileDetailsFragment,
  InvocationTestResultDetailsFragment,
  InvocationTestSummaryDetailsFragment,
} from "@/graphql/__generated__/graphql";
import { FILE_DETAILS_FRAGMENT } from "@/types/GraphqlFileFragment";
import { readableDurationFromMilliseconds } from "@/utils/time";
import {
  type TestSummaryDetailsTableRows,
  testActionOutputColumns,
  testResultColumns,
} from "./columns";

const TEST_LOG_FILE_NAME = "test.log";

interface Props {
  testSummary: InvocationTestSummaryDetailsFragment;
  testResults: InvocationTestResultDetailsFragment[];
}

export const TestSummaryPage: React.FC<Props> = ({
  testSummary,
  testResults,
}) => {
  const tableRows: TestSummaryDetailsTableRows[] = useMemo(() => {
    return testResults.flatMap(
      (tr) =>
        tr.testActionOutput?.map((file, index) => ({
          isFirstTestResultRow: index === 0,
          numberOfTestResultRows: tr.testActionOutput?.length ?? 0,
          id: tr.id,
          run: tr.run,
          shard: tr.shard,
          attempt: tr.attempt,
          status: tr.status ?? undefined,
          exitCode: tr.exitCode ?? undefined,
          strategy: tr.strategy ?? undefined,
          cacheLocation: cacheLocationFromBooleans(
            tr.cachedLocally,
            tr.cachedRemotely,
          ),
          testAttemptDurationInMs: tr.testAttemptDurationInMs ?? undefined,
          file: getFragmentData(FILE_DETAILS_FRAGMENT, file),
        })) ?? [],
    );
  }, [testResults]);

  const cacheLocation: CacheLocation = useMemo(
    () => cacheLocationFromTestResults(testResults),
    [testResults],
  );

  const logFile: FileDetailsFragment | undefined = useMemo(() => {
    if (testResults.length !== 1) {
      return undefined;
    }
    for (const fileFragment of testResults[0]?.testActionOutput ?? []) {
      const file = getFragmentData(FILE_DETAILS_FRAGMENT, fileFragment);
      if (file.filePath.path === TEST_LOG_FILE_NAME) {
        return file;
      }
    }
    return undefined;
  }, [testResults]);

  return (
    <PortalCard
      icon={<ExperimentFilled />}
      titleBits={["Test Summary"]}
      extraBits={[
        <LinkButton
          key="button-to-test-overview"
          buttonType="primary"
          to="/targets/$targetID/tests"
          params={{ targetID: testSummary.invocationTarget.target.id }}
        >
          Test Overview
        </LinkButton>,
      ]}
    >
      <Space direction="vertical" style={{ width: "100%" }} size="large">
        <Descriptions
          bordered
          layout="vertical"
          size="small"
          // TODO: Make this break columns based on the width of the content
          column={{ xs: 1, sm: 2, md: 3, lg: 3, xl: 6, xxl: 6 }}
          style={{ width: "max-content" }}
        >
          <Descriptions.Item label="Instance Name">
            {testSummary.invocationTarget.target.instanceName.name || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="Target Kind">
            {testSummary.invocationTarget.target.targetKind || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="Target Label">
            <Typography.Text copyable>
              {testSummary.invocationTarget.target.label || "-"}
            </Typography.Text>
          </Descriptions.Item>
          <Descriptions.Item label="Target Aspect">
            {testSummary.invocationTarget.target.aspect || "-"}
          </Descriptions.Item>

          <Descriptions.Item label="Status">
            <TestStatusTag key="status" status={testSummary.overallStatus} />
          </Descriptions.Item>
          <Descriptions.Item label="Cached">
            <CacheLocationTag cacheLocation={cacheLocation} />
          </Descriptions.Item>
        </Descriptions>

        <Space direction="horizontal" size="large">
          {testSummary.totalRunDurationInMs != null && (
            <Statistic
              title="Total duration"
              value={readableDurationFromMilliseconds(
                testSummary.totalRunDurationInMs,
                { smallestUnit: "ms" },
              )}
            />
          )}
          {testSummary.totalRunCount != null && (
            <Statistic
              title="Total run count"
              value={testSummary.totalRunCount}
            />
          )}
          {testSummary.runCount != null && (
            <Statistic title="Runs" value={testSummary.runCount} />
          )}
          {testSummary.attemptCount != null && (
            <Statistic title="Attempts" value={testSummary.attemptCount} />
          )}
          {testSummary.shardCount != null && (
            <Statistic title="Shards" value={testSummary.shardCount} />
          )}
          {testSummary.totalNumCached != null && (
            <Statistic title="Runs cached" value={testSummary.totalNumCached} />
          )}
        </Space>

        {logFile && (
          <CasGqlFileViewer
            file={logFile}
            title="Test log"
            fileName={TEST_LOG_FILE_NAME}
          />
        )}

        {testResults.length >= 0 && (
          <>
            <Typography.Title level={3}>
              {testResults.length > 1 ? "Test results" : "Test action outputs"}
            </Typography.Title>
            <Table
              columns={
                testResults.length > 1
                  ? [...testResultColumns, ...testActionOutputColumns]
                  : testActionOutputColumns
              }
              bordered={true}
              style={{ width: "100%" }}
              dataSource={tableRows}
              size="small"
              pagination={false}
            />
          </>
        )}
      </Space>
    </PortalCard>
  );
};
