import type { TableColumnsType } from "antd";
import {
  type CacheLocation,
  CacheLocationTag,
} from "@/components/CacheLocationTag";
import { TestStatusTag } from "@/components/TestStatusTag";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { readableFileSize } from "@/utils/filesize";
import { readableDurationFromMilliseconds } from "@/utils/time";
import { generateFileUrlFromGraphqlFile } from "@/utils/urlGenerator";

export type TestSummaryDetailsTableRows = {
  isFirstTestResultRow: boolean;
  numberOfTestResultRows: number;
  id: string;
  run: number;
  shard: number;
  attempt: number;
  status: string | undefined;
  exitCode: number | undefined;
  strategy: string | undefined;
  testAttemptDurationInMs: number | undefined;
  cacheLocation: CacheLocation;
  file: FileDetailsFragment;
};

const cellMergingLogic = (value: TestSummaryDetailsTableRows) => {
  if (value.isFirstTestResultRow) {
    return { rowSpan: value.numberOfTestResultRows };
  }
  return { rowSpan: 0 };
};

export const testResultColumns: TableColumnsType<TestSummaryDetailsTableRows> =
  [
    {
      key: "run",
      title: "Run",
      render: (_, record) => record.run,
      onCell: cellMergingLogic,
    },
    {
      key: "shard",
      title: "Shard",
      render: (_, record) => record.shard,
      onCell: cellMergingLogic,
    },
    {
      key: "attempt",
      title: "Attempt",
      render: (_, record) => record.attempt,
      onCell: cellMergingLogic,
    },
    {
      key: "strategy",
      title: "Strategy",
      render: (_, record) => record.strategy,
      onCell: cellMergingLogic,
    },
    {
      key: "status",
      title: "Status",
      render: (_, record) => <TestStatusTag status={record.status} />,
      onCell: cellMergingLogic,
    },
    {
      key: "exitCode",
      title: "Exit code",
      render: (_, record) => record.exitCode,
      onCell: cellMergingLogic,
    },
    {
      key: "cached",
      title: "Cached",
      render: (_, record) => (
        <CacheLocationTag cacheLocation={record.cacheLocation} />
      ),
      onCell: cellMergingLogic,
    },
    {
      key: "duration",
      title: "Duration",
      render: (_, record) =>
        readableDurationFromMilliseconds(record.testAttemptDurationInMs),
      onCell: cellMergingLogic,
    },
  ];

export const testActionOutputColumns: TableColumnsType<TestSummaryDetailsTableRows> =
  [
    {
      key: "fileSize",
      title: "File size",
      width: "min",
      render: (_, record) => readableFileSize(record.file.digest.sizeBytes),
    },
    {
      key: "file",
      title: "File",
      render: (_, record) => (
        <a href={generateFileUrlFromGraphqlFile(record.file)}>
          {record.file.filePath.path}
        </a>
      ),
    },
  ];
