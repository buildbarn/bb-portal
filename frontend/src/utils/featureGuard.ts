export class FeatureDisabledError extends Error {
  constructor() {
    super('FEATURE_DISABLED')
  }
}

export const requireFeature = (feature: Record<string, any> | undefined) => {
  return () => {
    if (feature === undefined) {
      throw new FeatureDisabledError();
    }
  }
}
