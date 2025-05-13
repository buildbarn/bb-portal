"use client";

import type { UrlObject } from "node:url";
import { useGrpcClients } from "@/context/GrpcClientsContext";
import {
  type Digest,
  type DigestFunction_Value,
  Directory,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { FileSystemAccessProfile } from "@/lib/grpc-client/buildbarn/fsac/fsac";
import type { FileSystemAccessProfileReference } from "@/lib/grpc-client/buildbarn/query/query";
import type { ByteStreamClient } from "@/lib/grpc-client/google/bytestream/bytestream";
import themeStyles from "@/theme/theme.module.css";
import {
  type BloomFilterReader,
  PathHashes,
  containsPathHashes,
  generateFileSystemReferenceQueryParams,
  readBloomFilter,
} from "@/utils/bloomFilter";
import { digestFunctionValueToString } from "@/utils/digestFunctionUtils";
import { fetchCasObjectAndParse } from "@/utils/fetchCasObject";
import { formatFileSizeFromString } from "@/utils/formatValues";
import { generateFileUrl } from "@/utils/urlGenerator";
import { DownOutlined, RightOutlined } from "@ant-design/icons";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, Flex, Space, Spin, Typography } from "antd";
import Link from "next/link";
import React, { useEffect } from "react";
import PortalAlert from "../PortalAlert";
import CopyBbClientdDirectoryButton from "./CopyBbClientdDirectoryButton";
import DownloadAsTarballButton from "./DownloadAsTarballButton";

const FETCH_STALE_TIME = 30000;

interface Params {
  instanceName: string;
  digestFunction: DigestFunction_Value;
  inputRootDigest: Digest;
  fileSystemAccessProfile: FileSystemAccessProfile | undefined;
  fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined;
}

const BrowserDirectory: React.FC<Params> = ({
  instanceName,
  digestFunction,
  inputRootDigest,
  fileSystemAccessProfile,
  fileSystemAccessProfileReference,
}) => {
  const bloomFilterReader = fileSystemAccessProfile
    ? readBloomFilter(fileSystemAccessProfile)
    : undefined;

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <RecursiveDirectoryNode
        instanceName={instanceName}
        digestFunction={digestFunction}
        directoryDigest={inputRootDigest}
        directoryName="Top level"
        isTopLevel={true}
        bloomFilterReader={bloomFilterReader}
        pathHashes={
          bloomFilterReader
            ? new PathHashes(
                fileSystemAccessProfileReference &&
                  BigInt(fileSystemAccessProfileReference?.pathHashesBaseHash),
              )
            : undefined
        }
        fileSystemAccessProfileRef={fileSystemAccessProfileReference}
      />
      <Space direction="vertical" size="small">
        {bloomFilterReader && (
          <Typography.Text>
            <strong>Note:</strong>{" "}
            <span className={themeStyles.colorSuccess}>Green</span> and{" "}
            <span className={themeStyles.colorFailure}>
              <s>red</s>
            </span>{" "}
            filenames above indicate which files and directories will be
            prefetched the next time a similar action executes. Though it is
            representative of what is actually accessed by the action, it may
            contain false positives and negatives.
          </Typography.Text>
        )}

        <Space direction="horizontal">
          <CopyBbClientdDirectoryButton
            instanceName={instanceName}
            digestFunction={digestFunction}
            inputRootDigest={inputRootDigest}
          />
          <DownloadAsTarballButton
            instanceName={instanceName}
            digestFunction={digestFunction}
            directoryDigest={inputRootDigest}
          />
        </Space>
      </Space>
    </Space>
  );
};

const fetchDirectory = async (
  casByteStreamClient: ByteStreamClient,
  instanceName: string,
  digestFunction: DigestFunction_Value,
  digest: Digest,
) => {
  return fetchCasObjectAndParse(
    casByteStreamClient,
    instanceName,
    digestFunction,
    digest,
    Directory,
  );
};

const RecursiveDirectoryNode: React.FC<{
  instanceName: string;
  digestFunction: DigestFunction_Value;
  directoryDigest: Digest;
  directoryName: string;
  isTopLevel: boolean;
  bloomFilterReader?: BloomFilterReader;
  pathHashes?: PathHashes;
  willBePrefetched?: boolean;
  fileSystemAccessProfileRef: FileSystemAccessProfileReference | undefined;
}> = ({
  instanceName,
  digestFunction,
  directoryDigest,
  directoryName,
  isTopLevel,
  bloomFilterReader,
  pathHashes,
  willBePrefetched,
  fileSystemAccessProfileRef,
}) => {
  const [expanded, setExpanded] = React.useState(isTopLevel);
  const queryClient = useQueryClient();
  const { casByteStreamClient } = useGrpcClients();

  const { data, isError, isPending, error } = useQuery({
    queryKey: [
      "browserDirectory",
      instanceName,
      digestFunction,
      directoryDigest,
    ],
    queryFn: fetchDirectory.bind(
      null,
      casByteStreamClient,
      instanceName,
      digestFunction,
      directoryDigest,
    ),
    staleTime: FETCH_STALE_TIME,
    refetchOnMount: "always",
  });

  // Prefetch all child directories. React-query will cache the results for us
  // and reuse them for the `useQuery` above.
  useEffect(() => {
    if (data) {
      for (const dirNode of data.directories) {
        if (dirNode.digest) {
          queryClient.prefetchQuery({
            queryKey: [
              "browserDirectory",
              instanceName,
              digestFunction,
              dirNode.digest,
            ],
            queryFn: fetchDirectory.bind(
              null,
              casByteStreamClient,
              instanceName,
              digestFunction,
              dirNode.digest,
            ),
            staleTime: FETCH_STALE_TIME,
          });
        }
      }
    }
  }, [casByteStreamClient, data, digestFunction, instanceName, queryClient]);

  const calcWillBePrefetched = (
    currentPathHashes: PathHashes | undefined = pathHashes,
  ) => {
    if (willBePrefetched === false) {
      return false;
    }
    if (bloomFilterReader === undefined || currentPathHashes === undefined) {
      return undefined;
    }
    return containsPathHashes(bloomFilterReader, currentPathHashes);
  };

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

  return (
    <>
      {!isTopLevel && (
        <DirectoryNode
          isDirectory={true}
          name={`${directoryName}/`}
          href={{
            pathname: `/browser/${instanceName}/blobs/${digestFunctionValueToString(
              digestFunction,
            )}/directory/${directoryDigest.hash}-${directoryDigest.sizeBytes}`,
            query: generateFileSystemReferenceQueryParams(
              fileSystemAccessProfileRef,
              pathHashes,
            ),
          }}
          sizeBytes={directoryDigest.sizeBytes}
          permissions="drwxr-xr-x"
          expanded={expanded}
          setExpanded={setExpanded}
          willBePrefetched={calcWillBePrefetched()}
        />
      )}

      {expanded && (
        <div style={{ marginLeft: isTopLevel ? "0px" : "32px" }}>
          {data.directories.map(
            (dirNode) =>
              dirNode.name &&
              dirNode.digest && (
                <RecursiveDirectoryNode
                  key={dirNode.digest.hash}
                  instanceName={instanceName}
                  digestFunction={digestFunction}
                  directoryDigest={dirNode.digest}
                  directoryName={dirNode.name}
                  isTopLevel={false}
                  bloomFilterReader={bloomFilterReader}
                  pathHashes={pathHashes?.appendComponent(dirNode.name)}
                  willBePrefetched={calcWillBePrefetched()}
                  fileSystemAccessProfileRef={fileSystemAccessProfileRef}
                />
              ),
          )}
          {data.files.map((file) => (
            <DirectoryNode
              key={file.name}
              name={file.name}
              href={{
                pathname: file.digest
                  ? generateFileUrl(
                      instanceName,
                      digestFunction,
                      file.digest,
                      file.name,
                    )
                  : undefined,
              }}
              sizeBytes={file.digest?.sizeBytes}
              permissions={`-r-${file.isExecutable ? "x" : "-"}r-${
                file.isExecutable ? "x" : "-"
              }r-${file.isExecutable ? "x" : "-"}`}
              willBePrefetched={calcWillBePrefetched(
                pathHashes?.appendComponent(file.name),
              )}
            />
          ))}
          {data.symlinks.map((symlink) => (
            <DirectoryNode
              key={symlink.name}
              name={`${symlink.name} -> ${symlink.target}`}
              permissions="lrwxrwxrwx"
            />
          ))}
        </div>
      )}
    </>
  );
};

const ROW_HEIGHT = 20;
const BUTTON_WIDTH = 32;
const BUTTON_PADDING = 8;

const DirectoryNode: React.FC<{
  isDirectory?: boolean;
  name: string;
  href?: UrlObject;
  sizeBytes?: string;
  permissions: string;
  expanded?: boolean;
  setExpanded?: (expanded: boolean) => void;
  willBePrefetched?: boolean;
}> = ({
  isDirectory = false,
  name,
  href,
  sizeBytes,
  permissions,
  expanded,
  setExpanded,
  willBePrefetched,
}) => {
  const indent = isDirectory ? "0px" : `${BUTTON_WIDTH + BUTTON_PADDING}px`;

  const formattedFileName = () => {
    switch (willBePrefetched) {
      case true:
        return <span className={themeStyles.colorSuccess}>{name}</span>;
      case false:
        return (
          <span className={themeStyles.colorFailure}>
            <s>{name}</s>
          </span>
        );
      case undefined:
        return <span>{name}</span>;
    }
  };

  return (
    <Flex
      justify="space-between"
      style={{ height: `${ROW_HEIGHT}px`, marginLeft: indent }}
    >
      <Flex align="center" gap={BUTTON_PADDING}>
        {isDirectory && expanded !== undefined && setExpanded !== undefined && (
          <Button
            type="text"
            onClick={() => setExpanded(!expanded)}
            style={{ width: `${BUTTON_WIDTH}px`, height: `${ROW_HEIGHT}px` }}
          >
            {expanded ? <DownOutlined /> : <RightOutlined />}
          </Button>
        )}
        {href ? (
          <Link href={href}>{formattedFileName()}</Link>
        ) : (
          <Typography.Text>{formattedFileName()}</Typography.Text>
        )}
      </Flex>
      <Space size="large">
        {sizeBytes && <pre>{formatFileSizeFromString(sizeBytes)}</pre>}
        <pre>{permissions}</pre>
      </Space>
    </Flex>
  );
};

export default BrowserDirectory;
