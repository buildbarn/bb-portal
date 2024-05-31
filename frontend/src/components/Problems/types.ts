import {BlobReference, ProblemInfoFragment, TargetProblem, TestProblem} from "@/graphql/__generated__/graphql";


export type Problem = ProblemInfoFragment | TargetProblem | TestProblem;
export type ActionOutputType = BlobReference;
