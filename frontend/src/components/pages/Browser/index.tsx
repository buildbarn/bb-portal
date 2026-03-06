
import BrowserActionGrid from "@/components/BrowserActionGrid";
import BrowserCommandGrid from "@/components/BrowserCommandGrid";
import BrowserDirectoryPage from "@/components/BrowserDirectoryPage";
import BrowserPreviousExecutionsPage from "@/components/BrowserPreviousExecutionsPage";
import Content from "@/components/Content";
import PortalCard from "@/components/PortalCard";
import { BrowserPageParams, BrowserPageType } from "@/types/BrowserPageType";
import {
  CalculatorOutlined,
  CalendarFilled,
  CodeOutlined,
  FolderOpenFilled,
  HistoryOutlined,
} from "@ant-design/icons";
import { Typography } from "antd";
import type React from "react";
import { BrowserSearchParams } from "@/routes/browser.$";

interface Params {
  params: BrowserPageParams
  search: BrowserSearchParams
}

export const BrowserPage: React.FC<Params> = ({ params, search }) => {
  const renderChild = () => {
    switch (params.browserPageType) {
      case BrowserPageType.Action:
        return (
          <PortalCard
            icon={<CalculatorOutlined />}
            titleBits={[<span key="title">Action</span>]}
          >
            <BrowserActionGrid browserPageParams={params} />
          </PortalCard>
        );

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
            <BrowserDirectoryPage browserPageParams={params} fileSystemAccessProfileReference={search.fileSystemAccessProfile} />
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
              <code>forceUploadTreesAndDirectories: true</code>{" "}
              in your <code>bb_worker.jsonnet</code>.
            </div>
          </PortalCard>
        );

      case BrowserPageType.HistoricalExecuteResponse:
        return (
          <PortalCard
            icon={<CalendarFilled />}
            titleBits={[<span key="title">Historical Execute Response</span>]}
          >
            <Typography.Title level={2}>
              Historical Execute Response
            </Typography.Title>
            <BrowserActionGrid
              browserPageParams={params}
              showTitle
            />
          </PortalCard>
        );

      case BrowserPageType.PreviousExecutionStats:
        return (
          <PortalCard
            icon={<HistoryOutlined />}
            titleBits={[<span key="title">Previous executions stats</span>]}
          >
            <BrowserPreviousExecutionsPage
              browserPageParams={params}
            />
          </PortalCard>
        );
      default:
        throw new Error(`Unknown browser page type: ${params.browserPageType}`);
    }
  };
  return <Content content={renderChild()} />;
};
