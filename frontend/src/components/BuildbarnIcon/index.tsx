import Icon from "@ant-design/icons";
import { theme } from "antd";
import BuildbarnIconSvg from "./BuildbarnIconSvg";

const ICON_SIZE = 20;

const { useToken } = theme;

const BuildbarnIcon: React.FC = () => {
  const { token } = useToken();

  const renderIconImage = () => {
    return (
      <Icon
        component={() => (
          <BuildbarnIconSvg
            width={ICON_SIZE}
            height={ICON_SIZE}
            stroke={token.colorText}
            name="Home"
          />
        )}
      />
    );
  };

  return <Icon component={renderIconImage} />;
};

export default BuildbarnIcon;
