import { Link } from "@tanstack/react-router";
import { Descriptions, Space, Typography } from "antd";
import type React from "react";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";
import type { FileSystemAccessProfileReference } from "@/types/FileSystemAccessProfileReference";
import { PATH_HASH_BASE_HASH } from "@/utils/bloomFilter";
import { readableFileSizeFromString } from "@/utils/filesize";
import { readableDurationFromProtobufDuration } from "@/utils/time";
import { generateBrowserSplat } from "@/utils/urlGenerator";
import { ActionProperties } from "../ActionProperties";
import { BrowserCommandDescription } from "../BrowserCommandDescription";
import CopyBbClientdCommandButton from "../BrowserCommandDescription/CopyBbClientdCommandButton";
import DownloadAsShellScriptButton from "../BrowserCommandDescription/DownloadAsShellScriptButton";
import { BrowserDirectory, type CParams } from "../BrowserDirectory";
import CopyBbClientdDirectoryButton from "../BrowserDirectory/CopyBbClientdDirectoryButton";
import DownloadAsTarballButton from "../BrowserDirectory/DownloadAsTarballButton";
import { DirectoryPrefetchDescription } from "../BrowserDirectory/directoryPrefetchDescription";
import BrowserPreviousExecutionsDisplay from "../BrowserPreviousExecutionsDisplay";
import BrowserResultDescription from "../BrowserResultDescription";
import ConditionalToolInvocationLink from "../ConditionalToolInvocationLink";
import ExecutionMetadataTimeline from "../ExecutionMetadataTimeline";
import FilesTable from "../FilesTable";
import {
  filesTableEntriesFromActionResultAndCommand,
  filesTableEntriesFromServerLogs,
} from "../FilesTable/utils";
import PortalAlert from "../PortalAlert";
import PropertyTagList from "../PropertyTagList";
import type { PropertyTagListEntry } from "../PropertyTagList/types";
import type { fetchBrowserActionGrid } from "./fetch";

interface Params {
  browserPageParams: BrowserPageParams;
  showTitle?: boolean;
  data: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  openDirsString: string | undefined;
}

const BrowserActionGrid: React.FC<Params> = ({
  browserPageParams,
  showTitle,
  data,
  openDirsString,
}) => {
  if (!data.action.commandDigest) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Error fetching action"
        description={
          "Error occurred while fetching data from the server; No command digest found."
        }
      />
    );
  }

  let fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined;

  if (data.fileSystemAccessProfile) {
    if (data.action.platform) {
      fileSystemAccessProfileReference = {
        digest: data.reducedActionDigest,
        pathHashesBaseHash: PATH_HASH_BASE_HASH,
      };
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
  const fileStructureData: CParams | undefined = data.action.inputRootDigest
    ? {
        instanceName: browserPageParams.instanceName,
        digestFunction: browserPageParams.digestFunction,
        digest: data.action.inputRootDigest,
        action: data.action,
        fileSystemAccessProfile: data.fileSystemAccessProfile,
        reducedActionDigest: data.reducedActionDigest,
        fileSystemAccessProfileReference,
      }
    : undefined;
  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <Space
        direction="vertical"
        size="middle"
        style={{
          display: "flex",
          minWidth: 0,
          overflow: "hidden",
        }}
      >
        {showTitle && (
          <Typography.Title level={2}>
            <Link
              to="/browser/$"
              params={{
                _splat: generateBrowserSplat(
                  browserPageParams.instanceName,
                  browserPageParams.digestFunction,
                  data.actionDigest,
                  BrowserPageType.Action,
                ),
              }}
              style={{ textDecoration: "underline" }}
            >
              Action
            </Link>
          </Typography.Title>
        )}
        <ActionProperties
          actionData={data}
          browserPageParams={browserPageParams}
        />
      </Space>

      {data.casCommand ? (
        <>
          <BrowserCommandDescription
            browserPageParams={browserPageParams}
            command={data.casCommand}
            commandDigest={data.action.commandDigest}
            showTitle={true}
          />
          {data.action.commandDigest && (
            <Space direction="horizontal">
              <CopyBbClientdCommandButton
                digestFunction={browserPageParams.digestFunction}
                instanceName={browserPageParams.instanceName}
                commandDigest={data.action.commandDigest}
              />
              <DownloadAsShellScriptButton
                digestFunction={browserPageParams.digestFunction}
                instanceName={browserPageParams.instanceName}
                commandDigest={data.action.commandDigest}
              />
            </Space>
          )}
        </>
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
              to="/browser/$"
              params={{
                _splat: generateBrowserSplat(
                  browserPageParams.instanceName,
                  browserPageParams.digestFunction,
                  data.action.inputRootDigest,
                  BrowserPageType.Directory,
                ),
              }}
              search={{
                fileSystemAccessProfile: fileSystemAccessProfileReference,
              }}
              style={{ textDecoration: "underline" }}
            >
              Input files
            </Link>
          </Typography.Title>
          {fileStructureData && (
            <BrowserDirectory
              baseData={fileStructureData}
              openDirsString={openDirsString}
              useBloomFilter={true}
            />
          )}
          <Space direction="vertical" size="small">
            <DirectoryPrefetchDescription
              prefetchDataExists={!!fileStructureData?.fileSystemAccessProfile}
            />
            {fileStructureData && (
              <Space direction="horizontal">
                <CopyBbClientdDirectoryButton
                  instanceName={fileStructureData.instanceName}
                  digestFunction={fileStructureData.digestFunction}
                  inputRootDigest={fileStructureData.digest}
                />
                <DownloadAsTarballButton
                  instanceName={fileStructureData.instanceName}
                  digestFunction={fileStructureData.digestFunction}
                  directoryDigest={fileStructureData.digest}
                />
              </Space>
            )}
          </Space>
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
          isPending={false}
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
              isPending={false}
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
                {readableDurationFromProtobufDuration(
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
              <ConditionalToolInvocationLink
                toolInvocationID={data.requestMetadata.toolInvocationId}
              />
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
                `${readableDurationFromProtobufDuration(data.posixResourceUsage.userTime)} user`}
              {data.posixResourceUsage.userTime &&
                data.posixResourceUsage.systemTime &&
                ","}{" "}
              {data.posixResourceUsage.systemTime &&
                `${readableDurationFromProtobufDuration(data.posixResourceUsage.systemTime)} system`}
            </Descriptions.Item>
            <Descriptions.Item label="Maximum resident set size">
              {readableFileSizeFromString(
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
              } files, having a total size of ${readableFileSizeFromString(
                data.filePoolResourceUsage.filesSizeBytesPeak,
              )}`}
            </Descriptions.Item>
            <Descriptions.Item label="Reads">
              {`${
                data.filePoolResourceUsage.readsCount
              } operations, having a total size of ${readableFileSizeFromString(
                data.filePoolResourceUsage.readsSizeBytes,
              )}`}
            </Descriptions.Item>
            <Descriptions.Item label="Writes">
              {`${
                data.filePoolResourceUsage.writesCount
              } operations, having a total size of ${readableFileSizeFromString(
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
      {data.previousExecutionStats && data.reducedActionDigest && (
        <BrowserPreviousExecutionsDisplay
          browserParams={browserPageParams}
          previousExecutionStats={data.previousExecutionStats}
          showTitle={true}
          reducedActionDigest={data.reducedActionDigest}
        />
      )}
    </Space>
  );
};

export default BrowserActionGrid;
