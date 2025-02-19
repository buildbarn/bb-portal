import { type TableColumnsType, Typography } from "antd";
import type { ColumnType } from "antd/lib/table";
import Link from "next/link";

export interface FilesTableEntry {
  mode: string | undefined;
  size: string | undefined;
  filename: string;
  href: string | undefined;
}

const modeColumn: ColumnType<FilesTableEntry> = {
  key: "mode",
  title: "Mode",
  dataIndex: "mode",
};

const sizeColumn: ColumnType<FilesTableEntry> = {
  key: "size",
  title: "Size",
  dataIndex: "size",
};

const filenameColumn: ColumnType<FilesTableEntry> = {
  key: "filename",
  title: "Filename",
  render: (_, record) => {
    if (record.href) {
      return <Link href={record.href}>{record.filename}</Link>;
    }
    return <Typography.Text>{record.filename}</Typography.Text>;
  },
};

const getColumns = (): TableColumnsType<FilesTableEntry> => {
  return [modeColumn, sizeColumn, filenameColumn];
};

export default getColumns;
