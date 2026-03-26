import type { Empty } from "@/lib/grpc-client/google/protobuf/empty";

export class FeatureDisabledError extends Error {
  constructor() {
    super('FEATURE_DISABLED')
  }
}

export const requireFeature = (feature: Empty | undefined) => {
  return () => {
    if (feature === undefined) {
      throw new FeatureDisabledError();
    }
  }
}
