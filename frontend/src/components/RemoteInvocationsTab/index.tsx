import type {
  CompletedActionDetailsFragment,
  PageInfo,
} from "@/graphql/__generated__/graphql";
import { PageCursorTable } from "../PageCursorTable";
import type { GetPaginationUpdateLinkType } from "../PageCursorTable/types";
import { columns } from "./Columns";

interface Props {
  actions: CompletedActionDetailsFragment[];
  pageSize: number;
  getPaginationUpdateLink: GetPaginationUpdateLinkType;
  pageInfo: PageInfo;
}

export const CompletedAction: React.FC<Props> = ({
  actions,
  getPaginationUpdateLink,
  pageSize,
  pageInfo,
}) => {
  if (actions.length === 0) {
    return <div>No completed actions found for this invocation.</div>;
  }

  return (
    <div className="remote-invocations">
      <PageCursorTable<CompletedActionDetailsFragment>
        columns={columns}
        size="small"
        rowKey="id"
        dataSource={actions}
        pageInfo={pageInfo}
        getPaginationUpdateLink={getPaginationUpdateLink}
        pageSize={pageSize}
      />
    </div>
  );
};
