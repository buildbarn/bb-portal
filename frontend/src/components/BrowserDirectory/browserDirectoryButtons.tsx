import {
  DiffFilled,
  DiffOutlined,
  DoubleRightOutlined,
  DownOutlined,
  RightOutlined,
} from "@ant-design/icons";
import { Button, Space, Tooltip } from "antd";

interface Params {
  expand: (depth: number, selective?: boolean) => void;
  collapse: (times: number) => void;
  closeAll: () => void;
  loading: boolean;
  top: boolean;
  compare: boolean;
}
export const BrowserDirectoryButtons: React.FC<Params> = ({
  expand,
  collapse,
  closeAll,
  loading: working,
  top,
  compare,
}: Params) => {
  return (
    <Space.Compact>
      <Tooltip title="Expand 1 level">
        <Button
          onClick={() => expand(1)}
          style={
            top
              ? {
                  borderBottomLeftRadius: 0,
                }
              : {
                  borderTopLeftRadius: 0,
                }
          }
          disabled={working}
        >
          <DownOutlined />
        </Button>
      </Tooltip>
      <Tooltip title="Collapse 1 level">
        <Button onClick={() => collapse(1)}>
          <RightOutlined />
        </Button>
      </Tooltip>
      <Tooltip title="Expand 5 levels">
        <Button onClick={() => expand(5)} disabled={working}>
          <DoubleRightOutlined rotate={90} />
        </Button>
      </Tooltip>
      <Tooltip title="Collapse all">
        <Button
          onClick={closeAll}
          style={
            !compare
              ? top
                ? {
                    borderBottomRightRadius: 0,
                  }
                : {
                    borderTopRightRadius: 0,
                  }
              : {}
          }
        >
          <DoubleRightOutlined />
        </Button>
      </Tooltip>
      {compare && (
        <>
          <Tooltip title="Expand 1 yellow directory level">
            <Button onClick={() => expand(1, true)} disabled={working}>
              <DiffOutlined />
            </Button>
          </Tooltip>
          <Tooltip title="Expand 5 yellow directory levels">
            <Button
              onClick={() => expand(5, true)}
              style={
                top
                  ? {
                      borderBottomRightRadius: 0,
                    }
                  : {
                      borderTopRightRadius: 0,
                    }
              }
              disabled={working}
            >
              <DiffFilled />
            </Button>
          </Tooltip>
        </>
      )}
    </Space.Compact>
  );
};
