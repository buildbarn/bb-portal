import type React from 'react';
import { Space, Typography, Upload } from 'antd';
import { FileAddTwoTone } from '@ant-design/icons';
import type { UploadProps } from "antd";

const { Dragger } = Upload;

interface Props {
  label: string;
  description: React.ReactNode;
  action: string;
}

const Uploader: React.FC<Props> = ({ label, description, action }) => {
  const handleChange: UploadProps["onChange"] = ({ file }) => {
    if (file.response && !file.error) {
      file.url = file.response?.Location ?? file.url;
    }
  };

  return (
    <Dragger
      name="file"
      action={action}
      onChange={handleChange}
      accept=".ndjson"
      multiple
    >
      <Space direction="vertical" size="small">
        <Typography.Title level={1}>
          <FileAddTwoTone />
        </Typography.Title>
        <Typography.Text>{label}</Typography.Text>
        {description}
      </Space>
    </Dragger>
  );
};

export default Uploader;
