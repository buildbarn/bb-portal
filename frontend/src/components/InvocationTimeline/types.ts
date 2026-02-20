import type { SVGProps } from "react";
import type { CartesianTickItem } from "recharts/types/util/types";
import type { CommandLineData } from "../CommandLine";

export interface TickProps extends SVGProps<SVGTextElement> {
  payload: CartesianTickItem;
}

export interface InvocationInfo {
  invocationId: string;
  timestamps: number[];
  exitCodeName: string | undefined;
  timeSinceLastConnectionMillis: number | undefined;
  command?: CommandLineData;
  workflow?: string | null;
  job?: string | null;
  action?: string | null;
}
