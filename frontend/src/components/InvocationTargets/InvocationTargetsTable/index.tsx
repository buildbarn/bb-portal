import { useQuery } from "@apollo/client/react";
import { Space } from "antd";
import type { FilterValue } from "antd/es/table/interface";
import React from "react";
import {
  buildSetIndex,
  type NamedSetNode,
  type OutputGroupNode,
} from "@/components/Artifacts/graph";
import { ARTIFACT_GRAPH_QUERY } from "@/components/Artifacts/index.graphql";
import TargetArtifactFiles from "@/components/Artifacts/TargetArtifactFiles";
import {
  CursorTable,
  getNewPaginationVariables,
} from "@/components/CursorTable";
import type { PaginationVariables } from "@/components/CursorTable/types";
import PortalAlert from "@/components/PortalAlert";
import type {
  GetInvocationTargetsForInvocationQuery,
  InvocationTargetAbortReason,
  InvocationTargetWhereInput,
  TargetMetrics,
} from "@/graphql/__generated__/graphql";
import styles from "@/theme/theme.module.css";
import { parseGraphqlEdgeList } from "@/utils/parseGraphqlEdgeList";
import { InvocationTargetsMetrics } from "../InvocationTargetsMetrics";
import { columns, type InvocationTargetsTableRowType } from "./Columns";
import { GET_INVOCATION_TARGETS_FOR_INVOCATION } from "./graphql";

// targetKey identifies an artifact-graph target by label + aspect so it
// can be matched to an invocation-target row.
const targetKey = (label: string, aspect: string | null | undefined): string =>
  `${label} ${aspect ?? ""}`;

interface Props {
  invocationId: string;
  targetMetrics: TargetMetrics | undefined;
}

export const InvocationTargetsTable: React.FC<Props> = ({
  invocationId,
  targetMetrics,
}) => {
  const [paginationVariables, setPaginationVariables] =
    React.useState<PaginationVariables>(getNewPaginationVariables());
  const [filterVariables, setFilterVariables] =
    React.useState<InvocationTargetWhereInput>({});

  const { data, error, loading } =
    useQuery<GetInvocationTargetsForInvocationQuery>(
      GET_INVOCATION_TARGETS_FOR_INVOCATION,
      {
        variables: {
          invocationID: invocationId,
          where: filterVariables,
          ...paginationVariables,
        },
        fetchPolicy: "cache-first",
        notifyOnNetworkStatusChange: true,
      },
    );

  // The per-file artifact graph (only populated when the invocation was
  // recorded at the basic_and_target_and_artifacts save level). It is
  // matched to target rows by label + aspect and rendered inline in the
  // expandable row.
  const { data: artifactData } = useQuery(ARTIFACT_GRAPH_QUERY, {
    variables: { id: invocationId },
    fetchPolicy: "cache-first",
  });

  const { sets, outputGroupsByTarget } = React.useMemo(() => {
    const graph = artifactData?.getBazelInvocation?.artifactGraph;
    const setIndex = buildSetIndex(
      (graph?.namedSets ?? []).map(
        (s): NamedSetNode => ({
          id: s.id,
          childSetIds: s.childSetIds,
          files: s.files,
        }),
      ),
    );
    const byTarget = new Map<string, OutputGroupNode[]>();
    for (const t of graph?.targets ?? []) {
      byTarget.set(targetKey(t.label, t.aspect), t.outputGroups);
    }
    return { sets: setIndex, outputGroupsByTarget: byTarget };
  }, [artifactData]);

  const onFilterChange = (filters: Record<string, FilterValue | null>) => {
    const newFilters: InvocationTargetWhereInput[] = [];
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.length > 0) {
        switch (key) {
          case "target-kind":
            newFilters.push({
              hasTargetWith: [{ targetKindContainsFold: value[0] as string }],
            });
            break;
          case "label":
            newFilters.push({
              hasTargetWith: [{ labelContainsFold: value[0] as string }],
            });
            break;
          case "aspect":
            newFilters.push({
              hasTargetWith: [{ aspectContainsFold: value[0] as string }],
            });
            break;
          case "success":
            if (value.length === 1) {
              newFilters.push({ success: value[0] as boolean });
            }
            break;
          case "abort-reason":
            newFilters.push({
              abortReasonIn: value as InvocationTargetAbortReason[],
            });
            break;
        }
      }
    });
    setFilterVariables({ and: newFilters });
  };

  if (error) {
    return (
      <PortalAlert
        type="error"
        message={
          error?.message ||
          "An unknown error occurred while fetching invocation targets."
        }
        showIcon
        className={styles.alert}
      />
    );
  }

  const invocationTargets = parseGraphqlEdgeList<InvocationTargetsTableRowType>(
    data?.getBazelInvocation?.invocationTargets,
  );

  return (
    <>
      <InvocationTargetsMetrics
        targetMetrics={targetMetrics}
        invocationTargetsCount={data?.getBazelInvocation?.numTotal.totalCount}
        invocationTargetsBuiltSuccessfully={
          data?.getBazelInvocation?.numSuccessful.totalCount
        }
        invocationTargetsSkipped={
          data?.getBazelInvocation?.numSkipped.totalCount
        }
      />
      <CursorTable<InvocationTargetsTableRowType>
        rowKey={"id"}
        size="small"
        columns={columns}
        loading={loading}
        dataSource={invocationTargets}
        expandable={{
          rowExpandable: (record) => {
            const hasFailure =
              record.failureMessage !== null &&
              record.failureMessage !== undefined &&
              record.failureMessage !== "";
            const hasArtifacts = outputGroupsByTarget.has(
              targetKey(record.target.label, record.target.aspect),
            );
            return hasFailure || hasArtifacts;
          },
          expandedRowRender: (record) => {
            const hasFailure =
              record.failureMessage !== null &&
              record.failureMessage !== undefined &&
              record.failureMessage !== "";
            const groups = outputGroupsByTarget.get(
              targetKey(record.target.label, record.target.aspect),
            );
            return (
              <Space
                direction="vertical"
                style={{ paddingLeft: "48px", display: "flex" }}
              >
                {hasFailure && (
                  <>
                    <strong>Failure Message:</strong>
                    <pre
                      style={{
                        whiteSpace: "pre-wrap",
                        overflowWrap: "break-word",
                      }}
                    >
                      {record.failureMessage}
                    </pre>
                  </>
                )}
                {groups && groups.length > 0 && (
                  <TargetArtifactFiles sets={sets} outputGroups={groups} />
                )}
              </Space>
            );
          },
        }}
        showSorterTooltip={{ target: "sorter-icon" }}
        onChange={(_pagination, filters, _sorter, _extra) =>
          onFilterChange(filters)
        }
        pagination={{
          position: "bottom",
          justify: "end",
          size: "middle",
        }}
        pageInfo={{
          startCursor:
            data?.getBazelInvocation?.invocationTargets.pageInfo.startCursor,
          endCursor:
            data?.getBazelInvocation?.invocationTargets.pageInfo.endCursor,
          hasNextPage:
            data?.getBazelInvocation?.invocationTargets.pageInfo.hasNextPage,
          hasPreviousPage:
            data?.getBazelInvocation?.invocationTargets.pageInfo
              .hasPreviousPage,
        }}
        paginationVariables={paginationVariables}
        setPaginationVariables={setPaginationVariables}
      />
    </>
  );
};
