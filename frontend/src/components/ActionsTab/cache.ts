export interface ActionCacheStatus {
  label: "Remote hit" | "Disk hit" | "Hit" | "No hit" | "Unknown";
  color: "success" | "default";
  description: string;
}

export const getActionCacheStatus = (
  cacheHit: boolean | null | undefined,
  runner: string | null | undefined,
): ActionCacheStatus => {
  if (cacheHit === null || cacheHit === undefined) {
    return {
      label: "Unknown",
      color: "default",
      description: "Cache status was not reported by Bazel",
    };
  }

  if (!cacheHit) {
    return {
      label: "No hit",
      color: "default",
      description: "Bazel reported no disk or remote cache hit",
    };
  }

  if (runner === "remote cache hit") {
    return {
      label: "Remote hit",
      color: "success",
      description: "Bazel reported a remote cache hit",
    };
  }
  if (runner === "disk cache hit") {
    return {
      label: "Disk hit",
      color: "success",
      description: "Bazel reported a disk cache hit",
    };
  }
  return {
    label: "Hit",
    color: "success",
    description: "Bazel reported a cache hit",
  };
};
