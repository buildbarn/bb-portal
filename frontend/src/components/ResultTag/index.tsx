import { type GlobalToken, theme } from "antd";
import type { CSSProperties } from "react";

export type TagVariables = {
  icon: React.ReactNode;
  color: string;
  text?: string;
};

interface Props {
  tagVars: TagVariables;
}

const { useToken } = theme;

const ResultTag: React.FC<Props> = ({ tagVars }) => {
  const { token } = useToken();
  const backgroundColorKey = `${tagVars.color}-1` as keyof GlobalToken;
  const textColorKey = `${tagVars.color}-7` as keyof GlobalToken;
  const backgroundColor = token[backgroundColorKey];
  const textColor = token[textColorKey];

  const style: CSSProperties = {
    borderRadius: 4,
    boxShadow: "rgba(0, 0, 0, 0.06) 0px 1px 2px 0px inset",
    display: "inline-block",
    fontSize: token.fontSizeSM,
    paddingBottom: "2px",
    paddingLeft: "7px",
    paddingRight: "7px",
    paddingTop: "2px",
  };

  if (typeof backgroundColor === "string" && typeof textColor === "string") {
    style.backgroundColor = backgroundColor;
    style.color = textColor;
  }

  return (
    <span style={style}>
      <span
        style={
          tagVars.text !== undefined && tagVars.text.length > 0
            ? { marginRight: "7px" }
            : undefined
        }
      >
        {tagVars.icon}
      </span>
      {tagVars.text}
    </span>
  );
};

export default ResultTag;
