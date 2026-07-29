import { theme } from "antd";
import type { CSSProperties } from "react";

const { useToken } = theme;

interface Props extends React.HTMLProps<HTMLSpanElement> {
  children: React.ReactNode;
  bordered?: boolean;
  style?: CSSProperties;
}

const CodeText: React.FC<Props> = ({ bordered, children, style, ...props }) => {
  const { token } = useToken();
  return (
    <span
      style={Object.assign(
        {
          fontFamily: token.fontFamilyCode,
          fontSize: token.fontSizeSM,
        },
        bordered === true && {
          backgroundColor: "rgba(150, 150, 150, 0.1)",
          borderColor: "rgba(100, 100, 100, 0.2)",
          borderStyle: "solid",
          borderWidth: 1,
          borderRadius: token.borderRadiusXS,
          padding: token.paddingXXS,
        },
        style,
      )}
      {...props}
    >
      {children}
    </span>
  );
};

export default CodeText;
