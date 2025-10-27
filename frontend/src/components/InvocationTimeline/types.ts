import type { SVGProps } from "react";
import type { CartesianTickItem } from "recharts/types/util/types";

export interface TickProps extends SVGProps<SVGTextElement> {
  payload: CartesianTickItem;
}

export interface InvocationInfo {
  invocationId: string;
  timestamps: (number | undefined)[];
  exitCode?: string;
}
