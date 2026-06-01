import { RadiusUprightOutlined } from "@ant-design/icons";
import { useQuery } from "@apollo/client/react";
import { Alert, Collapse, Input, Table, Typography } from "antd";
import type { TableColumnsType } from "antd/lib";
import type React from "react";
import { useMemo, useState } from "react";
import PortalCard from "../PortalCard";
import ArtifactFileTable from "./ArtifactFileTable";
import {
  type ArtifactGraphData,
  buildSetIndex,
  countSetFiles,
  type NamedSetNode,
  type TargetNode,
  walkSetFiles,
} from "./graph";
import { ARTIFACT_GRAPH_QUERY } from "./index.graphql";

const TARGET_PAGE_SIZE = 20;
const FILE_PAGE_SIZE = 50;

interface OutputGroupPanelProps {
  sets: ReadonlyMap<string, NamedSetNode>;
  rootSetIds: readonly string[];
}

const OutputGroupPanel: React.FC<OutputGroupPanelProps> = ({
  sets,
  rootSetIds,
}) => {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const files = useMemo(
    () => Array.from(walkSetFiles(sets, rootSetIds)),
    [sets, rootSetIds],
  );
  return (
    <ArtifactFileTable
      files={files}
      page={page}
      pageSize={FILE_PAGE_SIZE}
      onPageChange={(p) => setPage(p)}
      search={search}
      onSearchChange={(v) => {
        setSearch(v);
        setPage(1);
      }}
    />
  );
};

interface ExpandedTargetProps {
  sets: ReadonlyMap<string, NamedSetNode>;
  target: TargetNode;
}

const ExpandedTarget: React.FC<ExpandedTargetProps> = ({ sets, target }) => {
  // Common case: a single output group named "default" — render its files
  // directly without an extra expander.
  if (
    target.outputGroups.length === 1 &&
    target.outputGroups[0].name === "default"
  ) {
    return (
      <OutputGroupPanel
        sets={sets}
        rootSetIds={target.outputGroups[0].rootSetIds}
      />
    );
  }
  const items = target.outputGroups.map((group) => {
    const fileCount = countSetFiles(sets, group.rootSetIds);
    return {
      key: group.name,
      label: (
        <Typography.Text>
          {group.name}{" "}
          <Typography.Text type="secondary">
            ({fileCount} file{fileCount !== 1 ? "s" : ""}
            {group.incomplete ? ", incomplete" : ""})
          </Typography.Text>
        </Typography.Text>
      ),
      children: <OutputGroupPanel sets={sets} rootSetIds={group.rootSetIds} />,
    };
  });
  return <Collapse items={items} />;
};

interface TargetRow {
  key: string;
  label: string;
  target: TargetNode;
  outputGroupSummary: string;
}

interface Props {
  invocationId: string;
}

const ArtifactFilesCard: React.FC<Props> = ({ invocationId }) => {
  const { data, loading, error } = useQuery(ARTIFACT_GRAPH_QUERY, {
    variables: { id: invocationId },
  });

  const graph = useMemo<ArtifactGraphData | null>(() => {
    const ag = data?.getBazelInvocation?.artifactGraph;
    if (!ag) return null;
    return {
      namedSets: ag.namedSets.map((s) => ({
        id: s.id,
        childSetIds: s.childSetIds,
        files: s.files,
      })),
      targets: ag.targets.map((t) => ({
        label: t.label,
        aspect: t.aspect,
        outputGroups: t.outputGroups,
      })),
    };
  }, [data]);

  const sets = useMemo(() => buildSetIndex(graph?.namedSets ?? []), [graph]);

  const [page, setPage] = useState(1);
  const [targetSearch, setTargetSearch] = useState("");

  if (loading && !data) return null;
  if (!graph || graph.targets.length === 0) {
    // Suppress the card entirely when there is nothing to show.
    return null;
  }

  const filteredTargets = targetSearch
    ? graph.targets.filter((t) =>
        t.label.toLowerCase().includes(targetSearch.toLowerCase()),
      )
    : graph.targets;

  const rows: TargetRow[] = filteredTargets.map((t, i) => ({
    key: `${t.label}|${t.aspect ?? ""}|${i}`,
    label: t.label + (t.aspect ? ` (${t.aspect})` : ""),
    target: t,
    outputGroupSummary: t.outputGroups
      .map((g) => `${g.name}: ${countSetFiles(sets, g.rootSetIds)}`)
      .join(", "),
  }));

  const columns: TableColumnsType<TargetRow> = [
    {
      title: (
        <Input.Search
          allowClear
          placeholder="Filter by target label"
          value={targetSearch}
          onChange={(e) => {
            setTargetSearch(e.target.value);
            setPage(1);
          }}
          style={{ maxWidth: 400 }}
          size="small"
        />
      ),
      dataIndex: "label",
      key: "label",
      render: (label: string) => (
        <Typography.Text copyable={{ text: label }}>{label}</Typography.Text>
      ),
    },
    {
      title: "Output Groups",
      dataIndex: "outputGroupSummary",
      key: "outputGroupSummary",
      render: (s: string) => <Typography.Text>{s}</Typography.Text>,
    },
  ];

  return (
    <PortalCard
      type="inner"
      icon={<RadiusUprightOutlined />}
      titleBits={[<span key="label">Artifact Files</span>]}
    >
      {error && (
        <Alert
          type="error"
          message={error.message}
          style={{ marginBottom: 8 }}
        />
      )}
      <Table<TargetRow>
        rowKey="key"
        columns={columns}
        dataSource={rows.slice(
          (page - 1) * TARGET_PAGE_SIZE,
          page * TARGET_PAGE_SIZE,
        )}
        expandable={{
          expandedRowRender: (record) => (
            <ExpandedTarget sets={sets} target={record.target} />
          ),
        }}
        pagination={{
          current: page,
          pageSize: TARGET_PAGE_SIZE,
          total: rows.length,
          showSizeChanger: false,
          onChange: (p) => setPage(p),
        }}
      />
    </PortalCard>
  );
};

export default ArtifactFilesCard;
