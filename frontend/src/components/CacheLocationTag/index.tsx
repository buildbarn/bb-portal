import {
  CheckCircleFilled,
  CloseCircleFilled,
  QuestionCircleFilled,
} from "@ant-design/icons";
import type React from "react";
import ResultTag from "@/components/ResultTag";

export enum CacheLocation {
  Remote,
  Local,
  NotCached,
  Unknown,
}

export const cacheLocationFromBooleans = (
  cachedLocally: boolean | undefined | null,
  cachedRemotely: boolean | undefined | null,
): CacheLocation => {
  if (cachedLocally === true) {
    return CacheLocation.Local;
  }
  if (cachedRemotely === true) {
    return CacheLocation.Remote;
  }
  if (cachedLocally === false || cachedRemotely === false) {
    return CacheLocation.NotCached;
  }
  return CacheLocation.Unknown;
};

interface TestResultCacheLocation {
  cachedLocally?: boolean | null | undefined;
  cachedRemotely?: boolean | null | undefined;
}

export const cacheLocationFromTestResults = (
  testResults: TestResultCacheLocation[] | null | undefined,
): CacheLocation => {
  if (!testResults || testResults.length === 0) {
    return CacheLocation.Unknown;
  }
  if (testResults.every((tr) => tr.cachedRemotely === true)) {
    return CacheLocation.Remote;
  }
  if (testResults.every((tr) => tr.cachedLocally === true)) {
    return CacheLocation.Local;
  }
  if (testResults.some((tr) => tr.cachedLocally === false)) {
    return CacheLocation.NotCached;
  }
  if (testResults.some((tr) => tr.cachedRemotely === false)) {
    return CacheLocation.NotCached;
  }
  return CacheLocation.Unknown;
};

interface Props {
  cacheLocation: CacheLocation;
  hideText?: boolean;
}

export const CacheLocationTag: React.FC<Props> = ({
  cacheLocation,
  hideText,
}) => {
  let text: string;
  let color: string;
  let icon: React.ReactNode;

  switch (cacheLocation) {
    case CacheLocation.Remote:
      color = "green";
      text = "Remote";
      icon = <CheckCircleFilled />;
      break;
    case CacheLocation.Local:
      color = "blue";
      text = "Local";
      icon = <CheckCircleFilled />;
      break;
    case CacheLocation.NotCached:
      color = "red";
      text = "Not cached";
      icon = <CloseCircleFilled />;
      break;
    case CacheLocation.Unknown:
      color = "red";
      text = "Unknown";
      icon = <QuestionCircleFilled />;
      break;
  }

  return (
    <ResultTag
      tagVars={{
        color: color,
        icon: icon,
        text: hideText ? undefined : text,
      }}
    />
  );
};
