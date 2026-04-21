import { theme } from "antd";
import { diffWordsWithSpace } from "diff";
import { InnerBrowserCommandDescription } from "@/components/BrowserCommandDescription";
import type {
  Command,
  Command_EnvironmentVariable,
  Digest,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import { styleMap } from "../../BrowserDirectory/utils";

const { useToken } = theme;

interface BrowserCommandDescriptionParamsDelta {
  browserPageParams: BrowserPageParams;
  command: Command;
  commandDigest: Digest;
  showTitle: boolean;
  compareCommand: Command;
  mergedMode?: boolean; // For comparing actions
}
const BrowserCommandDescriptionDelta: React.FC<
  BrowserCommandDescriptionParamsDelta
> = ({
  browserPageParams,
  command,
  commandDigest,
  showTitle,
  compareCommand,
  mergedMode,
}: BrowserCommandDescriptionParamsDelta) => {
  const { token } = useToken();
  const mergedEnvironmentVariables: (Command_EnvironmentVariable & {
    style: React.CSSProperties;
  })[] = command.environmentVariables.map((envVar) => {
    const compareEnvVar = compareCommand.environmentVariables.find(
      (compareEnvVar) =>
        compareEnvVar.name === envVar.name &&
        compareEnvVar.value === envVar.value,
    );
    return {
      ...envVar,
      style: compareEnvVar ? {} : styleMap("unique", token),
    };
  });
  if (mergedMode) {
    compareCommand?.environmentVariables.forEach((compareEnvVar) => {
      const envVar = command.environmentVariables.find(
        (baseEnvVar) =>
          baseEnvVar.name === compareEnvVar.name &&
          baseEnvVar.value === compareEnvVar.value,
      );
      if (!envVar) {
        mergedEnvironmentVariables.push({
          ...compareEnvVar,
          style: styleMap("missing", token),
        });
      }
    });
  }

  const argumentStyles: React.CSSProperties[] = [];
  const comparedFirstArgument = diffWordsWithSpace(
    command.arguments[0],
    compareCommand.arguments[0],
  );
  const comparedArguments = diffWordsWithSpace(
    command.arguments.slice(1).join(" "),
    compareCommand.arguments.slice(1).join(" "),
  );
  const mergedArguments = [];

  for (let i = 0; i < comparedFirstArgument.length; i++) {
    if (!comparedFirstArgument[i].added && !comparedFirstArgument[i].removed) {
      mergedArguments.push(comparedFirstArgument[i].value);
      argumentStyles.push({ fontWeight: "bold" });
      continue;
    }
    if (comparedFirstArgument[i].added) {
      mergedArguments.push(comparedFirstArgument[i].value);
      argumentStyles.push({ ...styleMap("unique", token), fontWeight: "bold" });
    }
    if (comparedFirstArgument[i].removed && mergedMode) {
      mergedArguments.push(comparedFirstArgument[i].value);
      argumentStyles.push({
        ...styleMap("missing", token),
        fontWeight: "bold",
      });
    }
  }
  mergedArguments.push(" ");
  argumentStyles.push({});
  for (let i = 0; i < comparedArguments.length; i++) {
    if (!comparedArguments[i].added && !comparedArguments[i].removed) {
      mergedArguments.push(comparedArguments[i].value);
      argumentStyles.push({});
      continue;
    }
    if (comparedArguments[i].added) {
      mergedArguments.push(comparedArguments[i].value);
      argumentStyles.push(styleMap("unique", token));
    }
    if (comparedArguments[i].removed && mergedMode) {
      mergedArguments.push(comparedArguments[i].value);
      argumentStyles.push(styleMap("missing", token));
    }
  }

  const mergedWorkingDirectory =
    command.workingDirectory || compareCommand.workingDirectory ? (
      <>
        <p
          style={
            mergedMode &&
            command.workingDirectory !== compareCommand?.workingDirectory
              ? styleMap("unique", token)
              : {}
          }
        >
          {command.workingDirectory}
        </p>
        {mergedMode &&
          compareCommand &&
          command.workingDirectory !== compareCommand.workingDirectory && (
            <p style={styleMap("missing", token)}>
              {compareCommand.workingDirectory}
            </p>
          )}
      </>
    ) : undefined;
  return (
    <InnerBrowserCommandDescription
      browserPageParams={browserPageParams}
      command={{
        ...command,
        environmentVariables: mergedEnvironmentVariables,
        argumentStyles,
        workingDirectory: mergedWorkingDirectory,
        arguments: mergedArguments,
      }}
      commandDigest={commandDigest}
      showTitle={showTitle}
    />
  );
};

export { BrowserCommandDescriptionDelta };
