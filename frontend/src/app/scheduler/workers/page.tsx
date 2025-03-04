"use client";

import Content from "@/components/Content";
import PortalAlert from "@/components/PortalAlert";
import PortalCard from "@/components/PortalCard";
import WorkersGrid from "@/components/WorkersGrid";
import type {
  SizeClassQueueName,
  WorkerState,
} from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import type { ListWorkerFilterType } from "@/types/ListWorkerFilterType";
import { CalendarFilled } from "@ant-design/icons";
import { useSearchParams } from "next/navigation";
import type React from "react";
import { useMemo } from "react";

const Page: React.FC = () => {
  const searchParams = useSearchParams();

  const { listWorkerFilterType, sizeClassQueueName, paginationCursor } =
    useMemo(() => {
      const listWorkerFilterTypeParam = searchParams.get(
        "listWorkerFilterType",
      );
      const sizeClassQueueNameUrlParam = searchParams.get("sizeClassQueueName");
      const paginationCursorUrlParam = searchParams.get("paginationCursor");

      let decodedListWorkerFilterType: ListWorkerFilterType | null = null;
      let decodedSizeClassQueueName: SizeClassQueueName | undefined = undefined;
      let decodedPaginationCursor: WorkerState["id"] | undefined = undefined;

      try {
        if (listWorkerFilterTypeParam && sizeClassQueueNameUrlParam) {
          decodedListWorkerFilterType =
            listWorkerFilterTypeParam as ListWorkerFilterType;
          decodedSizeClassQueueName = JSON.parse(
            sizeClassQueueNameUrlParam,
          ) as SizeClassQueueName;
          if (paginationCursorUrlParam) {
            decodedPaginationCursor = JSON.parse(
              paginationCursorUrlParam,
            ) as WorkerState["id"];
          }
        }
      } catch (error) {
        console.error("Failed to decode URL parameters:", error);
      }

      return {
        listWorkerFilterType: decodedListWorkerFilterType,
        sizeClassQueueName: decodedSizeClassQueueName,
        paginationCursor: decodedPaginationCursor,
      };
    }, [searchParams]);

  return (
    <Content
      content={
        <PortalCard
          icon={<CalendarFilled />}
          titleBits={[<span key="title">Workers</span>]}
        >
          {!sizeClassQueueName || !listWorkerFilterType ? (
            <PortalAlert
              className="error"
              message="Error: Failed to decode URL worker identification."
            />
          ) : (
            <WorkersGrid
              listWorkerFilterType={listWorkerFilterType}
              sizeClassQueueName={sizeClassQueueName}
              paginationCursor={paginationCursor}
            />
          )}
        </PortalCard>
      }
    />
  );
};

export default Page;
