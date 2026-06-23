import { Alert, Input, Table, Typography } from "antd";
import type { TableColumnsType } from "antd/lib";
import type React from "react";
import { readableFileSize } from "@/utils/filesize";
import ArtifactDownloadAction from "./ArtifactDownloadAction";
import type { ArtifactFileNode } from "./graph";

interface Props {
  files: ArtifactFileNode[];
  loading?: boolean;
  error?: Error | null;
  onPageChange: (page: number, pageSize: number) => void;
  page: number;
  pageSize: number;
  search?: string;
  onSearchChange?: (value: string) => void;
}

const columns: TableColumnsType<ArtifactFileNode> = [
  {
    title: "Path",
    dataIndex: "name",
    ellipsis: true,
    render: (name: string) => (
      <Typography.Text copyable={{ text: name }}>{name}</Typography.Text>
    ),
  },
  {
    title: "Size",
    dataIndex: "sizeBytes",
    align: "right",
    render: (size?: number | null) =>
      size != null ? readableFileSize(size) : "—",
  },
  {
    title: "Digest",
    dataIndex: "digest",
    ellipsis: true,
    render: (d?: string | null) =>
      d ? (
        <Typography.Text type="secondary" copyable={{ text: d }}>
          {`${d.slice(0, 12)}…`}
        </Typography.Text>
      ) : (
        "—"
      ),
  },
  {
    title: "",
    key: "download",
    align: "right",
    render: (_, file) => (
      <ArtifactDownloadAction
        downloadUrl={file.downloadUrl}
        uri={file.uri}
        fileName={file.name}
      />
    ),
  },
];

const ArtifactFileTable: React.FC<Props> = ({
  files,
  loading,
  error,
  onPageChange,
  page,
  pageSize,
  search,
  onSearchChange,
}) => {
  const filtered = search
    ? files.filter((f) => f.name.toLowerCase().includes(search.toLowerCase()))
    : files;
  const start = (page - 1) * pageSize;
  const visible = filtered.slice(start, start + pageSize);
  return (
    <>
      {error && (
        <Alert
          type="error"
          message={error.message}
          style={{ marginBottom: 8 }}
        />
      )}
      {onSearchChange !== undefined && (
        <Input.Search
          allowClear
          placeholder="Filter by path"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          style={{ marginBottom: 8, maxWidth: 360 }}
        />
      )}
      <Table
        rowKey={(f) => `${f.name}|${f.digest ?? ""}`}
        columns={columns}
        dataSource={visible}
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total: filtered.length,
          showSizeChanger: false,
          onChange: onPageChange,
        }}
      />
    </>
  );
};

export default ArtifactFileTable;
