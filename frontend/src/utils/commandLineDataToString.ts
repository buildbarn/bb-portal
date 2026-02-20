import type {
  CommandLineData,
  CommandLineOptions,
} from "@/components/CommandLine";

export const commandLineDataToString = (cmd: CommandLineData) => {
  if (!cmd) return "unknown";
  return [
    cmd.executable,
    cmd.command,
    ...cmd.options.map(
      (x: CommandLineOptions) => `--${x?.option}${x?.value && `=${x?.value}`}`,
    ),
    ...cmd.residual,
  ].join(" ");
};
