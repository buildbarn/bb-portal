import type { SVGProps } from "react";
import type { CartesianTickItem } from "recharts/types/util/types";
import type { InvocationTag } from "@/graphql/__generated__/graphql";
import type { CommandLineData } from "../CommandLine";
import type { InvocationResult } from "../InvocationResultTag/enum";

export interface TickProps extends SVGProps<SVGTextElement> {
  payload: CartesianTickItem;
}

export interface InvocationInfo {
  invocationId: string;
  timestamps: number[];
  invocationStatus: InvocationResult;
  command?: CommandLineData;
  tags: Omit<InvocationTag, "bazelInvocation">[];
}
