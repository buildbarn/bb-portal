import themeStyles from "@/theme/theme.module.css";
import { BuildOutlined } from "@ant-design/icons";
import { Space, Table, Typography } from "antd";
import type React from "react";
import getColumns, { type FilesTableEntry } from "./Columns";

type Props = {
  entries: FilesTableEntry[];
  isPending: boolean;
};

const FilesTable: React.FC<Props> = ({ entries, isPending }) => {
  return (
    <Table
      columns={getColumns()}
      bordered={true}
      style={{ width: "100%" }}
      dataSource={entries}
      size="small"
      pagination={false}
      rowKey={(item) => item.filename}
      rowClassName={() => themeStyles.compactTable}
      locale={{
        emptyText: isPending ? (
          <Typography.Text
            disabled
            className={themeStyles.tableEmptyTextTypography}
          >
            <Space>
              <BuildOutlined />
              Loading...
            </Space>
          </Typography.Text>
        ) : (
          <Typography.Text
            disabled
            className={themeStyles.tableEmptyTextTypography}
          >
            <Space>
              <BuildOutlined />
              No files found
            </Space>
          </Typography.Text>
        ),
      }}
    />
  );
};

export default FilesTable;
