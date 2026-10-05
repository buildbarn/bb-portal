import { useQueries } from "@tanstack/react-query";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import type { CriticalPath } from "./parseCriticalPath";
import { getCriticalPath } from "./utils";

// Profiles are immutable, so a critical path derived from a profile digest
// never needs to be refetched within a session.
const CRITICAL_PATH_STALE_TIME = Number.POSITIVE_INFINITY;

const getCriticalPathQueryKey = (profile: FileDetailsFragment) => [
  "criticalPath",
  profile.digest.rev2InstanceName,
  profile.digest.digestFunction,
  profile.digest.hash,
  profile.digest.sizeBytes,
];

export const useCriticalPaths = (
  profiles: Array<FileDetailsFragment | null | undefined>,
): Array<{
  data: CriticalPath | undefined;
  isPending: boolean;
  isError: boolean;
  error: Error | null;
}> => {
  return useQueries({
    queries: profiles.map((profile) => {
      const enabled = profile != null;
      return {
        queryKey: enabled
          ? getCriticalPathQueryKey(profile)
          : ["criticalPath", "", "", "", 0],
        queryFn: () => {
          if (!enabled) {
            // Disabled query, this should never be invoked.
            throw new Error("Critical path profile is not available");
          }
          return getCriticalPath(profile);
        },
        enabled,
        staleTime: CRITICAL_PATH_STALE_TIME,
      };
    }),
    combine: (results) =>
      results.map((result) => ({
        data: result.data,
        isPending: result.isPending,
        isError: result.isError,
        error: result.error,
      })),
  });
};
