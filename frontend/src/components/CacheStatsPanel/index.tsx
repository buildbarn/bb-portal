import { useQuery } from "@tanstack/react-query";
import { Col, Row, Skeleton, Statistic, Typography } from "antd";
import type React from "react";
import PortalAlert from "@/components/PortalAlert";
import { env } from "@/utils/env";

interface CacheStats {
  casWrites: number;
  casReads: number;
  casBytesWritten: number;
  casBytesRead: number;
  acWrites: number;
  acHitRate: number | null;
}

async function fetchInstant(query: string): Promise<number | null> {
  const url = `/api/v1/prometheus/api/v1/query?query=${encodeURIComponent(query)}`;
  const res = await fetch(url);
  if (!res.ok)
    throw new Error(`Prometheus error: ${res.status} ${res.statusText}`);
  const json = await res.json();
  if (json.status !== "success")
    throw new Error(json.error ?? "Prometheus query failed");
  const result = json.data.result as Array<{ value: [number, string] }>;
  if (!result.length) return null;
  const v = parseFloat(result[0].value[1]);
  return Number.isFinite(v) ? v : null;
}

async function fetchCacheStats(): Promise<CacheStats> {
  const [
    casWrites,
    casReads,
    casBytesWritten,
    casBytesRead,
    acWrites,
    acHitRate,
  ] = await Promise.all([
    fetchInstant(
      'sum(buildbarn_blobstore_blob_access_operations_duration_seconds_count{storage_type="cas",operation="Put",grpc_code="OK"})',
    ),
    fetchInstant(
      'sum(buildbarn_blobstore_blob_access_operations_duration_seconds_count{storage_type="cas",operation="Get",grpc_code="OK"})',
    ),
    fetchInstant(
      'sum(buildbarn_blobstore_blob_access_operations_blob_size_bytes_sum{storage_type="cas",operation="Put"})',
    ),
    fetchInstant(
      'sum(buildbarn_blobstore_blob_access_operations_blob_size_bytes_sum{storage_type="cas",operation="Get"})',
    ),
    fetchInstant(
      'sum(buildbarn_blobstore_blob_access_operations_duration_seconds_count{storage_type="ac",operation="Put",grpc_code="OK"})',
    ),
    fetchInstant(
      'sum(buildbarn_blobstore_blob_access_operations_duration_seconds_count{storage_type="ac",operation="Get",grpc_code="OK"}) / sum(buildbarn_blobstore_blob_access_operations_duration_seconds_count{storage_type="ac",operation="Get"})',
    ),
  ]);
  return {
    casWrites: casWrites ?? 0,
    casReads: casReads ?? 0,
    casBytesWritten: casBytesWritten ?? 0,
    casBytesRead: casBytesRead ?? 0,
    acWrites: acWrites ?? 0,
    acHitRate,
  };
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1_048_576) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1_073_741_824) return `${(bytes / 1_048_576).toFixed(1)} MB`;
  return `${(bytes / 1_073_741_824).toFixed(2)} GB`;
}

export const CacheStatsPanel: React.FC = () => {
  const prometheusEnabled = Boolean(env.prometheusUrl);

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["cacheStats"],
    queryFn: fetchCacheStats,
    refetchInterval: 60_000,
    enabled: prometheusEnabled,
  });

  if (!prometheusEnabled) {
    return (
      <Typography.Text type="secondary">
        Cache statistics require Prometheus. Set <code>prometheus_url</code> in
        the bb-portal configuration.
      </Typography.Text>
    );
  }

  if (isLoading) {
    return <Skeleton active paragraph={{ rows: 2 }} />;
  }

  if (isError) {
    return (
      <PortalAlert
        showIcon
        type="error"
        message="Could not load cache statistics"
        description={(error as Error).message}
      />
    );
  }

  const hitRate = data?.acHitRate;
  const hitRateValue =
    hitRate !== null && hitRate !== undefined
      ? `${(hitRate * 100).toFixed(1)}`
      : "—";
  const hitRateColor =
    hitRate !== null && hitRate !== undefined
      ? hitRate >= 0.5
        ? "#52C41A"
        : "#FA8C16"
      : undefined;

  return (
    <Row gutter={[24, 16]}>
      <Col xs={12} sm={8} md={4}>
        <Statistic
          title="CAS Writes"
          value={(data?.casWrites ?? 0).toLocaleString()}
          suffix="ops"
        />
      </Col>
      <Col xs={12} sm={8} md={4}>
        <Statistic
          title="CAS Reads"
          value={(data?.casReads ?? 0).toLocaleString()}
          suffix="ops"
        />
      </Col>
      <Col xs={12} sm={8} md={4}>
        <Statistic
          title="CAS Written"
          value={formatBytes(data?.casBytesWritten ?? 0)}
        />
      </Col>
      <Col xs={12} sm={8} md={4}>
        <Statistic
          title="CAS Read"
          value={formatBytes(data?.casBytesRead ?? 0)}
        />
      </Col>
      <Col xs={12} sm={8} md={4}>
        <Statistic
          title="AC Writes"
          value={(data?.acWrites ?? 0).toLocaleString()}
          suffix="ops"
        />
      </Col>
      <Col xs={12} sm={8} md={4}>
        <Statistic
          title="AC Hit Rate"
          value={hitRateValue}
          suffix={hitRate !== null && hitRate !== undefined ? "%" : ""}
          valueStyle={hitRateColor ? { color: hitRateColor } : undefined}
        />
      </Col>
    </Row>
  );
};
