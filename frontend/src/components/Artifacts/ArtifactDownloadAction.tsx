import {
  DownloadOutlined,
  ExportOutlined,
  MinusOutlined,
} from "@ant-design/icons";
import { Button, Tooltip } from "antd";
import type React from "react";

interface Props {
  downloadUrl?: string | null;
  uri?: string | null;
  fileName: string;
}

const ArtifactDownloadAction: React.FC<Props> = ({
  downloadUrl,
  uri,
  fileName,
}) => {
  if (!downloadUrl) {
    return (
      <Tooltip
        title={uri ? `Not downloadable through bb-portal: ${uri}` : "No URI"}
      >
        <MinusOutlined />
      </Tooltip>
    );
  }
  const isExternal =
    downloadUrl.startsWith("http://") || downloadUrl.startsWith("https://");
  return (
    <Button
      type="link"
      size="small"
      icon={isExternal ? <ExportOutlined /> : <DownloadOutlined />}
      href={downloadUrl}
      target={isExternal ? "_blank" : "_self"}
      rel={isExternal ? "noopener noreferrer" : undefined}
      download={isExternal ? undefined : fileName}
    >
      {isExternal ? "External" : "Download"}
    </Button>
  );
};

export default ArtifactDownloadAction;
