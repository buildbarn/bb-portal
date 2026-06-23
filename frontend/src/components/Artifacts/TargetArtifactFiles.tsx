import { Collapse, Typography } from "antd";
import type React from "react";
import { useMemo, useState } from "react";
import ArtifactFileTable from "./ArtifactFileTable";
import {
  countSetFiles,
  type NamedSetNode,
  type OutputGroupNode,
  walkSetFiles,
} from "./graph";

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

interface Props {
  sets: ReadonlyMap<string, NamedSetNode>;
  outputGroups: OutputGroupNode[];
}

// TargetArtifactFiles renders the output files for a single target. The
// common case — a single "default" output group — renders the file table
// directly; otherwise each output group gets its own collapsible panel.
const TargetArtifactFiles: React.FC<Props> = ({ sets, outputGroups }) => {
  if (outputGroups.length === 0) {
    return null;
  }
  if (outputGroups.length === 1 && outputGroups[0].name === "default") {
    return (
      <OutputGroupPanel sets={sets} rootSetIds={outputGroups[0].rootSetIds} />
    );
  }
  const items = outputGroups.map((group) => {
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

export default TargetArtifactFiles;
