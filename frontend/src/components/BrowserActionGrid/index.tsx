"use client";

import { useGrpcClients } from "@/context/GrpcClientsContext";
import { FileSystemAccessProfileReference } from "@/lib/grpc-client/buildbarn/query/query";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import {
  PATH_HASH_BASE_HASH,
  generateFileSystemReferenceQueryParams,
} from "@/utils/bloomFilter";
import {
  digestFunctionValueToString,
  getReducedActionDigest_SHA256,
} from "@/utils/digestFunctionUtils";
import { formatDuration, formatFileSizeFromString } from "@/utils/formatValues";
import { useQuery } from "@tanstack/react-query";
import { Descriptions, Space, Spin, Typography } from "antd";
import Link from "next/link";
import type React from "react";
import BrowserCommandDescription from "../BrowserCommandDescription";
import BrowserDirectory from "../BrowserDirectory";
import BrowserPreviousExecutionsDisplay from "../BrowserPreviousExecutionsDisplay";
import BrowserResultDescription from "../BrowserResultDescription";
import ExecutionMetadataTimeline from "../ExecutionMetadataTimeline";
import FilesTable from "../FilesTable";
import {
  filesTableEntriesFromActionResultAndCommand,
  filesTableEntriesFromServerLogs,
} from "../FilesTable/utils";
import PortalAlert from "../PortalAlert";
import PropertyTagList from "../PropertyTagList";
import type { PropertyTagListEntry } from "../PropertyTagList/types";
import CopyBbClientdActionButton from "./CopyBbClientdActionButton";
import { fetchBrowserActionGrid } from "./fetch";

interface Params {
  browserPageParams: BrowserPageParams;
  showTitle?: boolean;
}

const BrowserActionGrid: React.FC<Params> = ({
  browserPageParams,
  showTitle,
}) => {
  const {
    actionCacheClient,
    casByteStreamClient,
    initialSizeClassCacheClient,
    fileSystemAccessCacheClient,
  } = useGrpcClients();

  const { data, isError, isPending, error } = useQuery({
    queryKey: ["browserActionGrid", browserPageParams],
    queryFn: fetchBrowserActionGrid.bind(
      window,
      browserPageParams,
      actionCacheClient,
      casByteStreamClient,
      initialSizeClassCacheClient,
      fileSystemAccessCacheClient,
    ),
  });

  let fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined = undefined;

  if (isError) {
    return (
      <PortalAlert
        className="error"
        message={
          <>
            <Typography.Text>
              There was a problem communicating with the backend server:
            </Typography.Text>
            <pre>{String(error)}</pre>
          </>
        }
      />
    );
  }

  if (isPending) {
    return <Spin />;
  }

  if (data.fileSystemAccessProfile) {
    if (data.action.commandDigest && data.action.platform) {
      fileSystemAccessProfileReference =
        FileSystemAccessProfileReference.create({
          digest: getReducedActionDigest_SHA256(
            data.action.commandDigest,
            data.action.platform,
          ),
          pathHashesBaseHash: PATH_HASH_BASE_HASH,
        });
    }
  }

  const workerPropertyList = (): PropertyTagListEntry[] => {
    const workerData = JSON.parse(
      data.executeResponse?.result?.executionMetadata?.worker || "{}",
    );
    return Object.keys(workerData).map(
      (key) => ({ name: key, value: workerData[key] }) as PropertyTagListEntry,
    );
  };

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      {data.action ? (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          {showTitle && (
            <Typography.Title level={2}>
              <Link
                href={`/browser/${
                  browserPageParams.instanceName
                }/blobs/${digestFunctionValueToString(
                  browserPageParams.digestFunction,
                )}/action/${data.actionDigest.hash}-${
                  data.actionDigest.sizeBytes
                }`}
                style={{ textDecoration: "underline" }}
              >
                Action
              </Link>
            </Typography.Title>
          )}
          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            {data.action.timeout && (
              <Descriptions.Item label="Timeout:">
                {formatDuration(data.action.timeout)}
              </Descriptions.Item>
            )}
            <Descriptions.Item label="Do not cache">
              {data.action.doNotCache ? "Yes" : "No"}
            </Descriptions.Item>
            {data.action.platform && (
              <Descriptions.Item label="Platform properties">
                <PropertyTagList
                  propertyList={data.action.platform.properties}
                />
              </Descriptions.Item>
            )}
          </Descriptions>
          {data.action.commandDigest && data.action.inputRootDigest && (
            <CopyBbClientdActionButton
              browserPageParams={browserPageParams}
              actionDigest={data.actionDigest}
              commandDigest={data.action.commandDigest}
              inputRootDigest={data.action.inputRootDigest}
            />
          )}
        </Space>
      ) : (
        <Typography.Text>This action could not be found.</Typography.Text>
      )}

      {data.casCommand ? (
        <BrowserCommandDescription
          browserPageParams={browserPageParams}
          command={data.casCommand}
          commandDigest={data.action.commandDigest}
          showTitle={true}
        />
      ) : (
        <Typography.Text>
          The command of this action could not be found.
        </Typography.Text>
      )}

      <Space direction="vertical" size="middle" style={{ width: "100%" }}>
        <Typography.Title level={2}>Result</Typography.Title>
        {data.executeResponse ? (
          <BrowserResultDescription
            browserPageParams={browserPageParams}
            executeResponse={data.executeResponse}
            posixResourceUsage={data.posixResourceUsage}
            consoleOutputs={data.consoleOutputs}
          />
        ) : (
          <Typography.Text>
            The action result of this action could not be found.
          </Typography.Text>
        )}
      </Space>

      {data.action.inputRootDigest && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={2}>
            <Link
              href={{
                pathname: `/browser/${
                  browserPageParams.instanceName
                }/blobs/${digestFunctionValueToString(
                  browserPageParams.digestFunction,
                )}/directory/${data.action.inputRootDigest.hash}-${
                  data.action.inputRootDigest.sizeBytes
                }`,
                query: generateFileSystemReferenceQueryParams(
                  fileSystemAccessProfileReference,
                ),
              }}
              style={{ textDecoration: "underline" }}
            >
              Input files
            </Link>
          </Typography.Title>
          <BrowserDirectory
            instanceName={browserPageParams.instanceName}
            digestFunction={browserPageParams.digestFunction}
            inputRootDigest={data.action.inputRootDigest}
            fileSystemAccessProfile={data.fileSystemAccessProfile}
            fileSystemAccessProfileReference={fileSystemAccessProfileReference}
          />
        </Space>
      )}

      <Space direction="vertical" size="middle" style={{ width: "100%" }}>
        <Typography.Title level={2}>Output files</Typography.Title>
        <FilesTable
          entries={filesTableEntriesFromActionResultAndCommand(
            data.executeResponse?.result,
            data.casCommand,
            browserPageParams.instanceName,
            browserPageParams.digestFunction,
          )}
          isPending={isPending}
        />
      </Space>

      {data.executeResponse?.serverLogs &&
        Object.keys(data.executeResponse.serverLogs).length !== 0 && (
          <Space direction="vertical" size="middle" style={{ width: "100%" }}>
            <Typography.Title level={2}>Server logs</Typography.Title>
            <FilesTable
              entries={filesTableEntriesFromServerLogs(
                data.executeResponse.serverLogs,
                browserPageParams.instanceName,
                browserPageParams.digestFunction,
              )}
              isPending={isPending}
            />
          </Space>
        )}

      {data.executeResponse?.result?.executionMetadata && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={2}>Execution metadata</Typography.Title>

          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            <Descriptions.Item label="Worker">
              <PropertyTagList propertyList={workerPropertyList()} />
            </Descriptions.Item>
            <Descriptions.Item label="Timeline">
              <ExecutionMetadataTimeline
                executionMetadata={
                  data.executeResponse.result.executionMetadata
                }
              />
            </Descriptions.Item>
            {data.executeResponse.result.executionMetadata
              .virtualExecutionDuration && (
              <Descriptions.Item label="Virtual execution duration">
                {formatDuration(
                  data.executeResponse.result.executionMetadata
                    .virtualExecutionDuration,
                )}
              </Descriptions.Item>
            )}
          </Descriptions>
        </Space>
      )}

      {data.authenticationMetadata && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={3}>Authentication metadata</Typography.Title>

          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            <Descriptions.Item label="Publicly displayable">
              <pre>
                {JSON.stringify(data.authenticationMetadata.public, null, 2)}
              </pre>
            </Descriptions.Item>
          </Descriptions>
        </Space>
      )}

      {data.requestMetadata && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={3}>Request metadata</Typography.Title>

          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            {data.requestMetadata.toolDetails && (
              <Descriptions.Item label="Tool">
                {`${data.requestMetadata.toolDetails.toolName} ${data.requestMetadata.toolDetails.toolVersion}`}
              </Descriptions.Item>
            )}
            <Descriptions.Item label="Tool invocation ID">
              {data.requestMetadata.toolInvocationId}
            </Descriptions.Item>
            <Descriptions.Item label="Correlated invocations ID">
              {data.requestMetadata.correlatedInvocationsId}
            </Descriptions.Item>
            <Descriptions.Item label="Target ID">
              {data.requestMetadata.targetId}
            </Descriptions.Item>
            <Descriptions.Item label="Action mnemonic">
              {data.requestMetadata.actionMnemonic}
            </Descriptions.Item>
            <Descriptions.Item label="Action ID">
              {data.requestMetadata.actionId}
            </Descriptions.Item>
            <Descriptions.Item label="Configuration ID">
              {data.requestMetadata.configurationId}
            </Descriptions.Item>
          </Descriptions>
        </Space>
      )}

      {data.posixResourceUsage && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={3}>POSIX resource usage</Typography.Title>

          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            <Descriptions.Item label="CPU time">
              {data.posixResourceUsage.userTime &&
                `${formatDuration(data.posixResourceUsage.userTime)} user`}
              {data.posixResourceUsage.userTime &&
                data.posixResourceUsage.systemTime &&
                ","}{" "}
              {data.posixResourceUsage.systemTime &&
                `${formatDuration(data.posixResourceUsage.systemTime)} system`}
            </Descriptions.Item>
            <Descriptions.Item label="Maximum resident set size">
              {formatFileSizeFromString(
                data.posixResourceUsage.maximumResidentSetSize,
              )}
            </Descriptions.Item>
            <Descriptions.Item label="Paging">
              {`${data.posixResourceUsage.pageReclaims} reclaims, ${data.posixResourceUsage.pageFaults} faults, ${data.posixResourceUsage.swaps} swaps`}
            </Descriptions.Item>
            <Descriptions.Item label="Block operations">
              {`${data.posixResourceUsage.blockInputOperations} inputs, ${data.posixResourceUsage.blockOutputOperations} outputs`}
            </Descriptions.Item>
            <Descriptions.Item label="Messages">
              {`${data.posixResourceUsage.messagesSent} sent, ${data.posixResourceUsage.messagesReceived} received`}
            </Descriptions.Item>
            <Descriptions.Item label="Signals">
              {`${data.posixResourceUsage.signalsReceived} received`}
            </Descriptions.Item>
            <Descriptions.Item label="Context switches">
              {`${data.posixResourceUsage.voluntaryContextSwitches} voluntary, ${data.posixResourceUsage.involuntaryContextSwitches} involuntary`}
            </Descriptions.Item>
          </Descriptions>
        </Space>
      )}

      {data.filePoolResourceUsage && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={3}>
            File pool resource usage
          </Typography.Title>

          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            <Descriptions.Item label="Files created">
              {data.filePoolResourceUsage.filesCreated}
            </Descriptions.Item>
            <Descriptions.Item label="Peak usage">
              {`${
                data.filePoolResourceUsage.filesCountPeak
              } files, having a total size of ${formatFileSizeFromString(
                data.filePoolResourceUsage.filesSizeBytesPeak,
              )}`}
            </Descriptions.Item>
            <Descriptions.Item label="Reads">
              {`${
                data.filePoolResourceUsage.readsCount
              } operations, having a total size of ${formatFileSizeFromString(
                data.filePoolResourceUsage.readsSizeBytes,
              )}`}
            </Descriptions.Item>
            <Descriptions.Item label="Writes">
              {`${
                data.filePoolResourceUsage.writesCount
              } operations, having a total size of ${formatFileSizeFromString(
                data.filePoolResourceUsage.writesSizeBytes,
              )}`}
            </Descriptions.Item>
            <Descriptions.Item label="Truncates">
              {`${data.filePoolResourceUsage.truncatesCount} operations`}
            </Descriptions.Item>
          </Descriptions>
        </Space>
      )}

      {data.inputRootResourceUsage && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={3}>
            Input root resource usage
          </Typography.Title>

          <Descriptions
            column={1}
            size="small"
            bordered
            styles={{ label: { width: "25%" }, content: { width: "75%" } }}
          >
            <Descriptions.Item label="Directories">
              {`${data.inputRootResourceUsage.directoriesResolved} resolved, ${data.inputRootResourceUsage.directoriesRead} read`}
            </Descriptions.Item>
            <Descriptions.Item label="Files">
              {`${data.inputRootResourceUsage.filesRead} read`}
            </Descriptions.Item>
          </Descriptions>
        </Space>
      )}

      {data.monetaryResourceUsage && (
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Typography.Title level={3}>Monetary resource usage</Typography.Title>

          <Descriptions column={1} bordered>
            {Object.entries(data.monetaryResourceUsage.expenses).map(
              ([key, value]) => (
                <Descriptions.Item key={key} label={key}>
                  {`${value.currency} ${value.cost}`}
                </Descriptions.Item>
              ),
            )}
          </Descriptions>
        </Space>
      )}
      {data.previousExecutionStats &&
        data.action.commandDigest &&
        data.action.platform && (
          <BrowserPreviousExecutionsDisplay
            browserParams={browserPageParams}
            previousExecutionStats={data.previousExecutionStats}
            showTitle={true}
            reducedActionDigest={getReducedActionDigest_SHA256(
              data.action.commandDigest,
              data.action.platform,
            )}
          />
        )}
    </Space>
  );
};

export default BrowserActionGrid;
