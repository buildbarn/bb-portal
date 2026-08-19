import { useMemo } from "react";
import type {
  ActionTimingMetrics,
  ActionWhereInput,
  BazelInvocationActionFragment,
} from "@/graphql/__generated__/graphql";
import { PageCursorTable } from "../PageCursorTable";
import type {
  GetPaginationUpdateLinkType,
  PageInfo,
} from "../PageCursorTable/types";
import { tableFiltersToGraphqlWhere } from "../PageCursorTable/utils";
import { ActionDetails } from "./action";
import { getColumns } from "./columns";
import { ActionsMetrics } from "./metrics";

interface Props {
  actions: BazelInvocationActionFragment[];
  actionTimingMetrics: Pick<
    ActionTimingMetrics,
    | "totalExpectedTimeInMs"
    | "timeSavedByCacheHitsInMs"
    | "totalActions"
    | "timedActions"
    | "cacheHitActions"
    | "timedCacheHitActions"
  >;
  actionMnemonics: string[];
  configurationMnemonics: string[];
  executionPhaseTimeInMs: number | null | undefined;
  pageSize: number;
  onFilterChange: (where: ActionWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

export const ActionsTab: React.FC<Props> = ({
  actions,
  actionTimingMetrics,
  actionMnemonics,
  configurationMnemonics,
  executionPhaseTimeInMs,
  pageSize,
  onFilterChange,
  getPaginationUpdateLink,
  pageInfo,
}) => {
  const columns = useMemo(
    () => getColumns(actionMnemonics, configurationMnemonics),
    [actionMnemonics, configurationMnemonics],
  );

  return (
    <>
      <ActionsMetrics
        actionTimingMetrics={actionTimingMetrics}
        executionPhaseTimeInMs={executionPhaseTimeInMs}
      />
      <PageCursorTable
        size="small"
        columns={columns}
        dataSource={actions}
        rowKey="id"
        expandable={{
          expandedRowRender: (action) => <ActionDetails action={action} />,
        }}
        onChange={(_pagination, filters) => {
          onFilterChange(tableFiltersToGraphqlWhere(columns, filters));
        }}
        getPaginationUpdateLink={getPaginationUpdateLink}
        pageInfo={pageInfo}
        pageSize={pageSize}
      />
    </>
  );
};
