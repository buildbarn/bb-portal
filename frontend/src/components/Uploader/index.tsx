import React, { useState } from 'react';
import { Space, Typography, Upload, UploadFile, UploadProps } from 'antd'
import { FileAddTwoTone } from '@ant-design/icons';

const { Dragger } = Upload;

interface Props {
  label: string;
  description: React.ReactNode;
  action: string;
}

const Uploader: React.FC<Props> = ({ label, description, action }) => {
  const [fileList, setFileList] = useState<UploadFile[]>();

  const handleChange: UploadProps['onChange'] = (info) => {
    let newFileList = [...info.fileList];

    newFileList = newFileList.map((file) => {
      if (file.response && !file.error) {
        file.url = file.response.Location;
      }
      return file;
    });

    setFileList(newFileList);
  };

  return (
    <Dragger name="file" action={action} onChange={handleChange} accept=".ndjson" multiple>
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
