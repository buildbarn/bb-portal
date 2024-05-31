import {BuildStepResultEnum} from "@/components/BuildStepResultTag";

enum ItemStatus {
  UNKNOWN = 0,
  SUCCESS = 1,
  FAILURE = 2,
  ABORTED = 3,
}

const BuildStepResultRank: Record<BuildStepResultEnum, ItemStatus> = {
  [BuildStepResultEnum.UNKNOWN]: ItemStatus.UNKNOWN,
  [BuildStepResultEnum.NOT_BUILT]: ItemStatus.UNKNOWN,
  [BuildStepResultEnum.UNSTABLE]: ItemStatus.UNKNOWN,
  [BuildStepResultEnum.SUCCESS]: ItemStatus.SUCCESS,
  [BuildStepResultEnum.ABORTED]: ItemStatus.ABORTED,
  [BuildStepResultEnum.INTERRUPTED]: ItemStatus.ABORTED,
  [BuildStepResultEnum.BUILD_FAILURE]: ItemStatus.FAILURE,
  [BuildStepResultEnum.PARSING_FAILURE]: ItemStatus.FAILURE,
  [BuildStepResultEnum.TESTS_FAILED]: ItemStatus.FAILURE,
};

const byResultRank = (status: BuildStepResultEnum): number => {
  return BuildStepResultRank[status];
};

export default byResultRank;
