import { BuildStepResultEnum } from "@/components/BuildStepResultTag";

export const getBarColor = (exitCode: string | undefined): string => {
  // Corresponds to the tag colors in
  // @/components/BuildStepResultTag
  switch (exitCode) {
    case BuildStepResultEnum.SUCCESS:
      return "green";
    case BuildStepResultEnum.UNSTABLE:
      return "orange";
    case BuildStepResultEnum.PARSING_FAILURE:
      return "red";
    case BuildStepResultEnum.BUILD_FAILURE:
      return "red";
    case BuildStepResultEnum.TESTS_FAILED:
      return "red";
    case BuildStepResultEnum.NOT_BUILT:
      return "purple";
    case BuildStepResultEnum.ABORTED:
      return "cyan";
    case BuildStepResultEnum.INTERRUPTED:
      return "cyan";
    case BuildStepResultEnum.UNKNOWN:
      return "grey";
    default:
      return "grey";
  }
};
