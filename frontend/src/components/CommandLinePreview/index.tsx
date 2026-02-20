import { Typography, theme } from "antd";
import { commandLineDataToString } from "@/utils/commandLineDataToString";
import type { CommandLineData } from "../CommandLine";

interface Props {
  command: CommandLineData;
  copyable?: boolean;
  codeBlockWrapper?: boolean;
}

const { useToken } = theme;

const CommandLinePreview: React.FC<Props> = ({
  command,
  copyable,
  codeBlockWrapper,
}) => {
  const { token } = useToken();
  const cmd = commandLineDataToString(command);
  // Replaces all hyphens with non-breaking hypens
  const displayCommand = cmd.replaceAll("-", "\u2011");

  return (
    <>
      {copyable && (
        <span>
          <Typography.Text copyable={{ text: cmd ?? "Copy" }} />
        </span>
      )}
      <Typography.Text title={cmd}>
        {(codeBlockWrapper === true && (
          <span
            style={{
              backgroundColor: "rgba(150, 150, 150, 0.1)",
              borderColor: "rgba(100, 100, 100, 0.2)",
              borderRadius: token.borderRadiusXS,
              borderStyle: "solid",
              borderWidth: 1,
              display: "block",
              fontFamily: token.fontFamilyCode,
              fontSize: token.fontSizeSM,
              maxWidth: "100%",
              padding: token.paddingXXS,
              textWrap: "wrap",
            }}
          >
            {displayCommand}
          </span>
        )) || <code>{displayCommand}</code>}
      </Typography.Text>
    </>
  );
};

export default CommandLinePreview;
