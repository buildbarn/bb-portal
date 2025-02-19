import type { PlatformQueueState } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

export interface PlatformQueueTableState extends PlatformQueueState {
  numberOfSizeClasses: number;
  isFirstSizeClass: boolean;
}
