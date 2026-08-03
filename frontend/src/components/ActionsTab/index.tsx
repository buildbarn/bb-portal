import { useMemo } from "react";
import type {
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

interface Props {
  actions: BazelInvocationActionFragment[];
  actionMnemonics: string[];
  configurationMnemonics: string[];
  pageSize: number;
  onFilterChange: (where: ActionWhereInput[]) => void;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

export const ActionsTab: React.FC<Props> = ({
  actions,
  actionMnemonics,
  configurationMnemonics,
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
  );
};
