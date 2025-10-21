export interface PageInfo {
  startCursor?: string;
  endCursor?: string;
  hasNextPage?: boolean;
  hasPreviousPage?: boolean;
}

export type PaginationVariables = {
  pageSize: number;
  after?: string;
  first?: number;
  before?: string;
  last?: number;
};
