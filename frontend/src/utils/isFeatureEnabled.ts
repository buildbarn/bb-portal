import { env } from "next-runtime-env";

export enum FeatureType {
  BES = "bes",
  SCHEDULER = "scheduler",
  BROWSER = "browser",
}

export const isFeatureEnabled = (featureType: FeatureType): boolean => {
  switch (featureType) {
    case FeatureType.BES:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES") === "true";
    case FeatureType.SCHEDULER:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_SCHEDULER") === "true";
    case FeatureType.BROWSER:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BROWSER") === "true";
  }
};
