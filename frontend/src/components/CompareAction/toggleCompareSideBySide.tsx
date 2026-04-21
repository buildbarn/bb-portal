import { MergeCellsOutlined, SplitCellsOutlined } from "@ant-design/icons";
import { Grid, Space, Switch } from "antd";

interface Props {
  prefersCompareSideBySide: boolean;
  setPrefersCompareSideBySide: React.Dispatch<React.SetStateAction<boolean>>;
}
const ToggleCompareSideBySide: React.FC<Props> = ({
  prefersCompareSideBySide,
  setPrefersCompareSideBySide,
}) => {
  const screens = Grid.useBreakpoint();

  return (
    <>
      {/*On screens smaller than XL, CompareSideBySide mode is always active, so we don't need to show the toggle */}
      {screens.xl && (
        <Space
          style={{
            cursor: "pointer",
            justifyContent: "center",
            alignItems: "center",
            marginRight: "32px",
            scale: "1.4",
          }}
        >
          <Switch
            defaultChecked={prefersCompareSideBySide}
            onChange={(e) => {
              setPrefersCompareSideBySide(e);
            }}
            checkedChildren={
              <SplitCellsOutlined
                style={{
                  scale: "1.4",
                }}
              />
            }
            unCheckedChildren={
              <MergeCellsOutlined
                style={{
                  scale: "1.4",
                }}
              />
            }
          />
        </Space>
      )}
    </>
  );
};

export default ToggleCompareSideBySide;
