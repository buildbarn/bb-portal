import { CalculatorOutlined, CalendarFilled } from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Spin, Typography } from "antd";
import type React from "react";
import { useState } from "react";
import { PortalCard } from "@/components/PortalCard";
import { actionCacheClient } from "@/grpc/actionCacheClient";
import { casByteStreamClient } from "@/grpc/casByteStreamClient";
import { fileSystemAccessCacheClient } from "@/grpc/fileSystemAccessCacheClient";
import { initialSizeClassCacheClient } from "@/grpc/initialSizeClassCacheClient";
import type { BrowserSearchParams } from "@/routes/browser.$";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";
import BrowserActionGrid from "../BrowserActionGrid";
import { fetchBrowserActionGrid } from "../BrowserActionGrid/fetch";
import { CompareActionButtons } from "../CompareAction/CompareTargetButton";
import CompareActionGrid from "../CompareAction/Page";
import ToggleCompareSideBySide from "../CompareAction/toggleCompareSideBySide";
import PortalAlert from "../PortalAlert";

interface Params {
  params: BrowserPageParams;
  search: BrowserSearchParams;
}

export const BrowserActionPages: React.FC<Params> = ({ params, search }) => {
  const [prefersCompareSideBySide, setPrefersCompareSideBySide] =
    useState(true);
  const actionQuery = useQuery({
    queryKey: ["browserActionGrid", params.digest.hash],
    queryFn: () => {
      return fetchBrowserActionGrid(
        params,
        actionCacheClient,
        casByteStreamClient,
        initialSizeClassCacheClient,
        fileSystemAccessCacheClient,
      );
    },
    staleTime: 5 * 60 * 1000,
  });
  const compareActionQuery = useQuery({
    queryKey: ["browserActionGrid", search.comparedAction?.digest.hash],
    enabled: !!search.comparedAction,
    queryFn: () => {
      if (!search.comparedAction) {
        return undefined;
      }
      return fetchBrowserActionGrid(
        { ...search.comparedAction, browserPageType: BrowserPageType.Action },
        actionCacheClient,
        casByteStreamClient,
        initialSizeClassCacheClient,
        fileSystemAccessCacheClient,
      );
    },
    staleTime: 5 * 60 * 1000,
  });

  if (
    actionQuery.isPending ||
    (compareActionQuery.isPending && search.comparedAction)
  ) {
    return <Spin />;
  }

  if (
    actionQuery.isError ||
    compareActionQuery.isError ||
    (!compareActionQuery.data && search.comparedAction)
  ) {
    return (
      <PortalAlert
        showIcon
        type="error"
        title="Error fetching action"
        description={
          actionQuery.error?.message ||
          compareActionQuery.error?.message ||
          "Unknown error occurred while fetching data from the server."
        }
      />
    );
  }

  switch (params.browserPageType) {
    case BrowserPageType.Action:
      return (
        <PortalCard
          icon={<CalculatorOutlined />}
          titleBits={[
            search.comparedAction ? (
              <span key="title">Compare Actions</span>
            ) : (
              <span key="title">Action</span>
            ),
          ]}
          extraBits={[
            search.comparedAction !== undefined && (
              <ToggleCompareSideBySide
                prefersCompareSideBySide={prefersCompareSideBySide}
                setPrefersCompareSideBySide={setPrefersCompareSideBySide}
                key="toggleCompareSideBySide"
              />
            ),
            <CompareActionButtons
              params={params}
              comparing={search.comparedAction !== undefined}
              key="compareActionButtons"
            />,
          ]}
        >
          {search.comparedAction && compareActionQuery.data ? (
            <CompareActionGrid
              browserPageParams={params}
              compareParams={{
                ...search.comparedAction,
                browserPageType: BrowserPageType.Action,
              }}
              prefersCompareSideBySide={prefersCompareSideBySide}
              baseData={actionQuery.data}
              compareData={compareActionQuery.data}
              openDirsString={search.openDirs}
            />
          ) : (
            <BrowserActionGrid
              browserPageParams={params}
              data={actionQuery.data}
              openDirsString={search.openDirs}
            />
          )}
        </PortalCard>
      );
    case BrowserPageType.HistoricalExecuteResponse:
      return (
        <PortalCard
          icon={<CalendarFilled />}
          titleBits={[<span key="title">Historical Execute Response</span>]}
          extraBits={[
            <CompareActionButtons
              params={{
                ...params,
                browserPageType: BrowserPageType.Action,
                digest: actionQuery.data.actionDigest,
              }}
              comparing={search.comparedAction !== undefined}
              key="compareActionButtons"
            />,
          ]}
        >
          <Typography.Title level={2}>
            Historical Execute Response
          </Typography.Title>
          <BrowserActionGrid
            browserPageParams={params}
            data={actionQuery.data}
            showTitle
            openDirsString={search.openDirs}
          />
        </PortalCard>
      );
    default:
      throw new Error(`Unknown browser page type: ${params.browserPageType}`);
  }
};
