import {
  CodeOutlined,
  FolderOpenFilled,
  HistoryOutlined,
} from "@ant-design/icons";
import { Typography } from "antd";
import type React from "react";
import { BrowserActionPages } from "@/components/BrowserActionPages";
import BrowserCommandGrid from "@/components/BrowserCommandGrid";
import BrowserDirectoryPage from "@/components/BrowserDirectoryPage";
import BrowserPreviousExecutionsPage from "@/components/BrowserPreviousExecutionsPage";
import PortalCard from "@/components/PortalCard";
import type { BrowserSearchParams } from "@/routes/browser.$";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";

interface Params {
  params: BrowserPageParams;
  search: BrowserSearchParams;
}

export const BrowserPage: React.FC<Params> = ({ params, search }) => {
  switch (params.browserPageType) {
    case BrowserPageType.Action: // Intentional fall-though
    case BrowserPageType.HistoricalExecuteResponse:
      return <BrowserActionPages params={params} search={search} />;

    case BrowserPageType.Command:
      return (
        <PortalCard
          icon={<CodeOutlined />}
          titleBits={[<span key="title">Command</span>]}
        >
          <BrowserCommandGrid browserPageParams={params} />
        </PortalCard>
      );

    case BrowserPageType.Directory:
      return (
        <PortalCard
          icon={<FolderOpenFilled />}
          titleBits={[<span key="title">Directory</span>]}
        >
          <Typography.Title level={2}>Directory contents</Typography.Title>
          <BrowserDirectoryPage
            browserPageParams={params}
            fileSystemAccessProfileReference={search.fileSystemAccessProfile}
          />
        </PortalCard>
      );

    case BrowserPageType.Tree:
      return (
        <PortalCard
          icon={<FolderOpenFilled />}
          titleBits={[<span key="title">Tree directory</span>]}
        >
          <div>
            Tree objects are not supported. Please set{" "}
            <code>forceUploadTreesAndDirectories: true</code> in your{" "}
            <code>bb_worker.jsonnet</code>.
          </div>
        </PortalCard>
      );

    case BrowserPageType.PreviousExecutionStats:
      return (
        <PortalCard
          icon={<HistoryOutlined />}
          titleBits={[<span key="title">Previous executions stats</span>]}
        >
          <BrowserPreviousExecutionsPage browserPageParams={params} />
        </PortalCard>
      );
    default:
      throw new Error(`Unknown browser page type: ${params.browserPageType}`);
  }
};
