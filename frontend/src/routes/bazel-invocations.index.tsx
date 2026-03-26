import { createFileRoute } from "@tanstack/react-router";
import { BazelInvocationsPage } from "@/components/pages/BazelInvocations";
import { generatePageTitle } from "@/utils/generatePageTitle";

export const Route = createFileRoute("/bazel-invocations/")({
  component: BazelInvocationsPage,
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Invocations"]) }] }),
});
