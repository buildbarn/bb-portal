import { SwapOutlined } from "@ant-design/icons";
import { Link } from "@tanstack/react-router";
import {
  Button,
  Col,
  type GlobalToken,
  Grid,
  Row,
  Space,
  Typography,
  theme,
} from "antd";
import type React from "react";
import type { fetchBrowserActionGrid } from "@/components/BrowserActionGrid/fetch";
import { DirectoryPrefetchDescription } from "@/components/BrowserDirectory/directoryPrefetchDescription";
import FilesTable from "@/components/FilesTable";
import type { FilesTableEntry } from "@/components/FilesTable/Columns";
import { filesTableEntriesFromActionResultAndCommand } from "@/components/FilesTable/utils";
import {
  type BrowserPageParams,
  BrowserPageSchema,
} from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";
import { PATH_HASH_BASE_HASH } from "@/utils/bloomFilter";
import { generateBrowserSplat } from "@/utils/urlGenerator";
import { ActionPropertiesDelta } from "../../ActionProperties";
import { BrowserDirectory, type CParams } from "../../BrowserDirectory";
import type { CompareState } from "../../BrowserDirectory/types";
import { styleMap } from "../../BrowserDirectory/utils";
import { BrowserCommandDescriptionDelta } from "../BrowserCommandDescriptionDelta";
import styles from "./index.module.css";

const { useToken } = theme;

interface Params {
  browserPageParams: BrowserPageParams;
  compareParams: BrowserPageParams;
  prefersCompareSideBySide: boolean;
  baseData: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  compareData: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  openDirsString: string | undefined;
}

const CompareActionGrid: React.FC<Params> = ({
  browserPageParams,
  compareParams,
  prefersCompareSideBySide,
  baseData,
  compareData,
  openDirsString,
}) => {
  const { token } = useToken();

  const screens = Grid.useBreakpoint();
  const compareSideBySide = !!screens.xl && prefersCompareSideBySide;

  if (!baseData.action.inputRootDigest || !compareData.action.inputRootDigest) {
    return null;
  }

  const fileStructureData: CParams = {
    instanceName: browserPageParams.instanceName,
    digestFunction: browserPageParams.digestFunction,
    digest: baseData.action.inputRootDigest,
    action: baseData.action,
    fileSystemAccessProfile: baseData.fileSystemAccessProfile,
    reducedActionDigest: baseData.reducedActionDigest,
    fileSystemAccessProfileReference: {
      digest: baseData.reducedActionDigest,
      pathHashesBaseHash: PATH_HASH_BASE_HASH,
    },
  };
  const compareFileStructureData: CParams = {
    instanceName: compareParams.instanceName,
    digestFunction: compareParams.digestFunction,
    digest: compareData.action.inputRootDigest,
    action: compareData.action,
    fileSystemAccessProfile: compareData.fileSystemAccessProfile,
    reducedActionDigest: compareData.reducedActionDigest,
    fileSystemAccessProfileReference: {
      digest: compareData.reducedActionDigest,
      pathHashesBaseHash: PATH_HASH_BASE_HASH,
    },
  };
  const [outputFiles, compareOutputFiles, allFileData] = generateOutputFileData(
    baseData,
    browserPageParams,
    compareData,
    { ...compareParams, browserPageType: BrowserPageType.Action },
    token,
  );

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      <div
        className={[styles.container, screens.lg ? styles.large : ""].join(" ")}
      >
        <div className={[styles.left, styles.actionLabel].join(" ")}>
          <DigestTitle
            actionPageParams={browserPageParams}
            screens={screens}
            compareSideBySide={compareSideBySide}
            isCompare={false}
            token={token}
          />
        </div>
        <div className={[styles.right, styles.actionLabel].join(" ")}>
          <DigestTitle
            actionPageParams={compareParams}
            screens={screens}
            compareSideBySide={compareSideBySide}
            isCompare={true}
            token={token}
          />
        </div>
        <div className={styles.button}>
          <Link
            to="/browser/$"
            params={{
              _splat: generateBrowserSplat(
                compareParams.instanceName,
                compareParams.digestFunction,
                compareParams.digest,
                BrowserPageType.Action,
              ),
            }}
            search={{
              comparedAction: BrowserPageSchema.parse({
                ...browserPageParams,
                browserPageType: BrowserPageType.Action,
              }),
            }}
          >
            <Button
              type="primary"
              shape={!screens.lg ? "default" : "circle"}
              icon={<SwapOutlined rotate={!screens.lg ? 90 : 0} />}
            />
          </Link>
        </div>
      </div>
      <Row>
        {compareSideBySide ? (
          <Space direction="vertical">
            <p>
              <span style={styleMap("unique", token)}>Green</span> highlighted
              objects are unique to the action in that column.
            </p>
            <p>
              <span style={styleMap("missing", token)}>Red</span> highlighted
              objects exists only for the action in the opposite column.
            </p>
            <p>
              <span style={styleMap("different", token)}>Yellow</span>{" "}
              highlighted objects have the same name or path, but different
              values in the differnet actions. In the input folder, this means
              the folder has different objects inside it or a subfolder, or that
              a symlink has the same name but points at different locations.
            </p>

            <DirectoryPrefetchDescription prefetchDataExists={true} />
          </Space>
        ) : (
          <Space direction="vertical">
            <p>
              <span style={styleMap("unique", token)}>Green</span> highlighted
              objects are unique to the first action, which has its hash
              highlighted in{" "}
              <span style={styleMap("unique", token)}>Green</span> above.
            </p>
            <p>
              <span style={styleMap("missing", token)}>Red</span> highlighted
              objects are unique to the second action, which has its hash
              highlighted in <span style={styleMap("missing", token)}>Red</span>{" "}
              above.
            </p>
            <p>
              <span style={styleMap("different", token)}>Yellow</span>{" "}
              highlighted objects have the same name or path, but different
              values in the differnet actions. In the input folder, this means
              the folder has different objects inside it or a subfolder, or that
              a symlink has the same name but points at different locations.
            </p>
          </Space>
        )}
      </Row>
      <Row gutter={24}>
        <Col span={compareSideBySide ? 12 : 24}>
          <Space
            direction="vertical"
            size="middle"
            style={{
              display: "flex",
              minWidth: 0,
              overflow: "hidden",
            }}
          >
            <ActionPropertiesDelta
              actionData={baseData}
              compareActionData={compareData}
              browserPageParams={browserPageParams}
              mergedMode={!compareSideBySide}
            />
          </Space>
        </Col>
        {compareSideBySide && (
          <Col span={12}>
            <Space
              direction="vertical"
              size="middle"
              style={{
                display: "flex",
                minWidth: 0,
                overflow: "hidden",
              }}
            >
              <ActionPropertiesDelta
                actionData={compareData}
                compareActionData={baseData}
                browserPageParams={compareParams}
                mergedMode={!compareSideBySide}
              />
            </Space>
          </Col>
        )}
      </Row>
      <Row gutter={24}>
        <Col span={compareSideBySide ? 12 : 24}>
          {baseData.casCommand &&
          baseData.action.commandDigest &&
          compareData.casCommand ? (
            <BrowserCommandDescriptionDelta
              browserPageParams={browserPageParams}
              command={baseData.casCommand}
              commandDigest={baseData.action.commandDigest}
              showTitle={true}
              compareCommand={compareData.casCommand}
              mergedMode={!compareSideBySide}
            />
          ) : (
            <Typography.Text>
              Error finding command of this action.
            </Typography.Text>
          )}
        </Col>
        {compareSideBySide && (
          <Col span={12}>
            {compareData.casCommand &&
            compareData.action.commandDigest &&
            baseData.casCommand ? (
              <BrowserCommandDescriptionDelta
                browserPageParams={compareParams}
                command={compareData.casCommand}
                commandDigest={compareData.action.commandDigest}
                showTitle={true}
                compareCommand={baseData.casCommand}
                mergedMode={!compareSideBySide}
              />
            ) : (
              <Typography.Text>
                The command of this action could not be found.
              </Typography.Text>
            )}
          </Col>
        )}
      </Row>
      <Row gutter={24}>
        <Col span={24} style={{ marginTop: 12, marginBottom: 12 }}>
          <BrowserDirectory
            baseData={fileStructureData}
            compareData={compareFileStructureData}
            mergedMode={!compareSideBySide}
            openDirsString={openDirsString}
            useBloomFilter={compareSideBySide}
          />
        </Col>
      </Row>
      <Row gutter={24}>
        <Col span={compareSideBySide ? 12 : 24}>
          <Space direction="vertical" size="middle" style={{ width: "100%" }}>
            <Typography.Title level={2}>Output files</Typography.Title>
            <FilesTable
              entries={compareSideBySide ? outputFiles : allFileData}
              isPending={false}
            />
          </Space>
        </Col>
        {compareSideBySide && (
          <Col span={12}>
            <Space direction="vertical" size="middle" style={{ width: "100%" }}>
              <Typography.Title level={2}>Output files</Typography.Title>
              <FilesTable entries={compareOutputFiles} isPending={false} />
            </Space>
          </Col>
        )}
      </Row>
    </Space>
  );
};

const getStyledEntry = (
  entry: FilesTableEntry,
  type: CompareState,
  token: GlobalToken,
) => ({
  ...entry,
  style: {
    color: token.colorText,
    padding: "4px",
    borderRadius: "5px",
    ...styleMap(type, token),
  },
});

const generateOutputFileData = (
  actionData: Awaited<ReturnType<typeof fetchBrowserActionGrid>> | undefined,
  browserPageParams: BrowserPageParams,
  compareData: Awaited<ReturnType<typeof fetchBrowserActionGrid>> | undefined,
  compareParams: BrowserPageParams,
  token: GlobalToken,
) => {
  const outputFiles = filesTableEntriesFromActionResultAndCommand(
    actionData?.executeResponse?.result,
    actionData?.casCommand,
    browserPageParams.instanceName,
    browserPageParams.digestFunction,
  );

  const compareFiles = filesTableEntriesFromActionResultAndCommand(
    compareData?.executeResponse?.result,
    compareData?.casCommand,
    compareParams.instanceName,
    compareParams.digestFunction,
  );

  const compareMap = new Map(compareFiles.map((f) => [f.filename, f]));
  const outputMap = new Map(outputFiles.map((f) => [f.filename, f]));

  // 1. Process Output Files & Identify Differences
  const processedOutput = outputFiles.map((entry) => {
    const match = compareMap.get(entry.filename);

    if (!match) return getStyledEntry(entry, "unique", token);

    const isDifferent =
      entry.href !== match.href ||
      entry.mode !== match.mode ||
      entry.size !== match.size;
    return isDifferent ? getStyledEntry(entry, "different", token) : entry;
  });

  // 2. Process Compare Files
  const processedCompare = compareFiles.map((entry) => {
    const match = outputMap.get(entry.filename);

    if (!match) return getStyledEntry(entry, "unique", token);

    const isDifferent =
      entry.href !== match.href ||
      entry.mode !== match.mode ||
      entry.size !== match.size;
    return isDifferent ? getStyledEntry(entry, "different", token) : entry;
  });

  // 3. Build Merged Data for merged mode
  // We take the processed base output, and add "missing" entries for things only in compare
  const missingFromOutput = compareFiles
    .filter((f) => !outputMap.has(f.filename))
    .map((f) => getStyledEntry(f, "missing", token));

  const allFileData = [...processedOutput, ...missingFromOutput].sort((a, b) =>
    a.filename.localeCompare(b.filename),
  );

  return [processedOutput, processedCompare, allFileData];
};

type DigestTitleProps = {
  actionPageParams: BrowserPageParams;
  screens: Record<string, boolean>;
  token: GlobalToken;
  isCompare: boolean;
  compareSideBySide: boolean;
};

const DigestTitle = ({
  actionPageParams,
  screens,
  token,
  isCompare,
  compareSideBySide,
}: DigestTitleProps) => {
  return (
    <Typography.Title
      level={screens.xxl ? 4 : 5}
      style={{
        borderRadius: "8px",
        padding: "2px",
        ...styleMap("link", token),
        ...(!compareSideBySide &&
          styleMap(isCompare ? "missing" : "unique", token)),
      }}
    >
      <Link
        to="/browser/$"
        params={{
          _splat: generateBrowserSplat(
            actionPageParams.instanceName,
            actionPageParams.digestFunction,
            actionPageParams.digest,
            BrowserPageType.Action,
          ),
        }}
        style={{ textDecoration: "underline" }}
      >
        {actionPageParams.digest.hash}-{actionPageParams.digest.sizeBytes}
      </Link>
    </Typography.Title>
  );
};

export default CompareActionGrid;
