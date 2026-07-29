import { commandLineDataToString } from "@/utils/commandLineDataToString";
import CodeText from "../CodeText";
import type { CommandLineData } from "../CommandLine";
import CopyableIcon from "../CopyableIcon";

interface Props {
  command: CommandLineData;
  copyable?: boolean;
  codeBlockWrapper?: boolean;
}

const CommandLinePreview: React.FC<Props> = ({
  command,
  copyable,
  codeBlockWrapper,
}) => {
  const cmd = commandLineDataToString(command);
  // Replaces all hyphens with non-breaking hypens
  const displayCommand = cmd.replaceAll("-", "\u2011");

  return (
    <>
      {copyable && <CopyableIcon text={cmd} />}
      <CodeText
        title={cmd}
        bordered
        style={
          codeBlockWrapper
            ? {
                display: "block",
                textWrap: "wrap",
              }
            : undefined
        }
      >
        {displayCommand}
      </CodeText>
    </>
  );
};

export default CommandLinePreview;
