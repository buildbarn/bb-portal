import type { LinkOptions } from "@tanstack/react-router";
import { z } from "zod";
import type { BazelInvocationWhereInput } from "@/graphql/__generated__/graphql";

export interface PageInfo {
  startCursor?: string;
  endCursor?: string;
  hasNextPage?: boolean;
  hasPreviousPage?: boolean;
}

export const TablePaginationVarsSchema = z.object({
  pageSize: z.number().int().positive().optional(),
  pagination: z
    .object({
      after: z.string().optional(),
      first: z.number().int().positive().optional(),
      before: z.string().optional(),
      last: z.number().int().positive().optional(),
    })
    .optional(),
});

export type TablePaginationVars = z.infer<typeof TablePaginationVarsSchema>;

export type GetPaginationUpdateLinkType = (
  vars: TablePaginationVars,
) => LinkOptions;

export type OnBazelInvocationFilterChange = (
  where: BazelInvocationWhereInput[],
) => void;
