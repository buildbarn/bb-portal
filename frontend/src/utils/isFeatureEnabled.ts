import { env } from "next-runtime-env";

export enum FeatureType {
  BEP_UPLOAD = "bep_upload",
  BES = "bes",
  BES_PAGE_BUILDS = "bes_page_builds",
  BES_PAGE_INVOCATIONS = "bes_page_invitations",
  BES_PAGE_TARGETS = "bes_page_targets",
  BES_PAGE_TESTS = "bes_page_tests",
  BES_PAGE_TRENDS = "bes_page_trends",
  BROWSER = "browser",
  SCHEDULER = "scheduler",
  OPERATIONS = "operations",
}

export const isFeatureEnabled = (featureType: FeatureType): boolean => {
  switch (featureType) {
    case FeatureType.BEP_UPLOAD:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BEP_UPLOAD") === "true";
    case FeatureType.BES:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES") === "true";
    case FeatureType.BES_PAGE_BUILDS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_BUILDS") === "true";
    case FeatureType.BES_PAGE_INVOCATIONS:
      return (
        env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_INVOCATIONS") === "true"
      );
    case FeatureType.BES_PAGE_TRENDS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TRENDS") === "true";
    case FeatureType.BES_PAGE_TARGETS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TARGETS") === "true";
    case FeatureType.BES_PAGE_TESTS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BES_PAGE_TESTS") === "true";
    case FeatureType.SCHEDULER:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_SCHEDULER") === "true";
    case FeatureType.OPERATIONS:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_OPERATIONS") === "true";
    case FeatureType.BROWSER:
      return env("NEXT_PUBLIC_ENABLED_FEATURES_BROWSER") === "true";
  }
};
