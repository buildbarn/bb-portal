import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";

export interface FilesTableEntry {
  mode: string | undefined;
  size: string | undefined;
  filename: string;
  href: string | undefined;
}

const modeColumn: ColumnType<FilesTableEntry> = {
  key: "mode",
  title: "Mode",
  render: (_, record) => record.mode,
};

const sizeColumn: ColumnType<FilesTableEntry> = {
  key: "size",
  title: "Size",
  render: (_, record) => record.size,
};

const filenameColumn: ColumnType<FilesTableEntry> = {
  key: "filename",
  title: "Filename",
  render: (_, record) => {
    if (record.href) {
      return <a href={record.href}>{record.filename}</a>;
    }
    return <Typography.Text>{record.filename}</Typography.Text>;
  },
};

const getColumns = (): TableColumnsType<FilesTableEntry> => {
  return [modeColumn, sizeColumn, filenameColumn];
};

export default getColumns;
