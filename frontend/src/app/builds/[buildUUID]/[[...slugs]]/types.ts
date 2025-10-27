import type { FindBuildByUuidQuery } from "@/graphql/__generated__/graphql";

export type FindBuildFromUuidFragment = NonNullable<
  NonNullable<FindBuildByUuidQuery["getBuild"]>["invocations"]
>[number];
