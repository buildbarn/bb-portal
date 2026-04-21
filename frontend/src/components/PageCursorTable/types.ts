import { z } from "zod";

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
    .partial()
    .optional(),
});

export type TablePaginationVars = z.infer<typeof TablePaginationVarsSchema>;
export type OnTablePaginationChange = (vars: TablePaginationVars) => void;
