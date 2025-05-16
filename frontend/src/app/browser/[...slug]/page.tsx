"use client";

import BrowserActionGrid from "@/components/BrowserActionGrid";
import BrowserCommandGrid from "@/components/BrowserCommandGrid";
import BrowserDirectoryPage from "@/components/BrowserDirectoryPage";
import BrowserPreviousExecutionsPage from "@/components/BrowserPreviousExecutionsPage";
import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import { BrowserPageType } from "@/types/BrowserPageType";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { parseBrowserPageSlug } from "@/utils/parseBrowserPageSlug";
import {
  CalculatorOutlined,
  CalendarFilled,
  CodeOutlined,
  FolderOpenFilled,
  HistoryOutlined,
} from "@ant-design/icons";
import { Typography } from "antd";
import { notFound } from "next/navigation";
import type React from "react";

interface PageParams {
  params: {
    slug: Array<string>;
  };
}

const Page: React.FC<PageParams> = ({ params }) => {
  if (!isFeatureEnabled(FeatureType.BROWSER)) {
    return <PageDisabled />;
  }
  const browserPageParams = parseBrowserPageSlug(params.slug);

  if (browserPageParams === undefined) {
    notFound();
  }

  const renderChild = () => {
    switch (browserPageParams.browserPageType) {
      case BrowserPageType.Action:
        return (
          <PortalCard
            icon={<CalculatorOutlined />}
            titleBits={[<span key="title">Action</span>]}
          >
            <BrowserActionGrid browserPageParams={browserPageParams} />
          </PortalCard>
        );

      case BrowserPageType.Command:
        return (
          <PortalCard
            icon={<CodeOutlined />}
            titleBits={[<span key="title">Command</span>]}
          >
            <BrowserCommandGrid browserPageParams={browserPageParams} />
          </PortalCard>
        );

      case BrowserPageType.Directory:
        return (
          <PortalCard
            icon={<FolderOpenFilled />}
            titleBits={[<span key="title">Directory</span>]}
          >
            <Typography.Title level={2}>Directory contents</Typography.Title>
            <BrowserDirectoryPage browserPageParams={browserPageParams} />
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
              browserPageParams={browserPageParams}
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
              browserPageParams={browserPageParams}
            />
          </PortalCard>
        );
      default:
        return notFound();
    }
  };

  return <Content content={renderChild()} />;
};

export default Page;
