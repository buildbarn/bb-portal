import { env } from "next-runtime-env";

export enum FeatureType {
  BES = "bes",
  BEP_UPLOAD = "bep_upload",
  SCHEDULER = "scheduler",
  BROWSER = "browser",
  BES_PAGE_TESTS = "bes_page_tests",
  BES_PAGE_TARGETS = "bes_page_targets",
}

export const isFeatureEnabled = (featureType: FeatureType): boolean => {
  switch (featureType) {
    case FeatureType.BES:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES") === "true";
    case FeatureType.BEP_UPLOAD:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BEP_UPLOAD") === "true";
    case FeatureType.SCHEDULER:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_SCHEDULER") === "true";
    case FeatureType.BROWSER:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BROWSER") === "true";
    case FeatureType.BES_PAGE_TESTS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TESTS") === "true";
    case FeatureType.BES_PAGE_TARGETS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TARGETS") === "true";
  }
};
