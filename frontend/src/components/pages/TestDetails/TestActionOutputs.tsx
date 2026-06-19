import { Space, Tooltip } from "antd";
import type React from "react";
import type { GetTestDetailsQuery } from "@/graphql/__generated__/graphql";

type TestResult = NonNullable<
  NonNullable<
    NonNullable<
      NonNullable<GetTestDetailsQuery["findTestSummaries"]["edges"]>[number]
    >["node"]
  >["testResults"]
>[number];

/** Per-output reference as stored in the test_action_outputs JSON column. */
interface TestActionOutputEntry {
  name?: string;
  uri?: string;
  digest?: string;
  length?: number;
}

interface Props {
  results: ReadonlyArray<TestResult>;
}

/**
 * Render a compact list of download links for each test action output.
 *
 * Surfaces test.log, test.xml, and test.outputs/outputs.zip URIs directly
 * from the BES TestResult event. Users can click through to the underlying
 * blob (bytestream:// URIs are resolved server-side via the existing
 * /blob proxy; file:// URIs are passed through and let the browser handle
 * them, which it will refuse for security — correct outcome).
 *
 * Falls back to "—" when no outputs are present (older rows that predate
 * the schema migration, or test attempts that produced no output files).
 */
export const TestActionOutputs: React.FC<Props> = ({ results }) => {
  const outputs = collectUniqueOutputs(results);
  if (outputs.length === 0) {
    return <span style={{ color: "rgba(0,0,0,0.45)" }}>—</span>;
  }
  return (
    <Space size={[4, 0]} wrap>
      {outputs.map((o) => (
        <Tooltip
          key={`${o.name ?? "?"}:${o.uri ?? "?"}`}
          title={o.length ? `${o.length.toLocaleString()} bytes` : undefined}
        >
          <a
            href={uriToDownloadHref(o.uri ?? "")}
            rel="noopener noreferrer"
            target="_blank"
          >
            {shortName(o.name ?? "")}
          </a>
        </Tooltip>
      ))}
    </Space>
  );
};

function collectUniqueOutputs(
  results: ReadonlyArray<TestResult>,
): TestActionOutputEntry[] {
  const seen = new Map<string, TestActionOutputEntry>();
  for (const r of results) {
    const raw = (r as { testActionOutputs?: unknown }).testActionOutputs;
    if (!raw || typeof raw !== "object") continue;
    for (const value of Object.values(raw as Record<string, unknown>)) {
      const entry = value as TestActionOutputEntry;
      const key = entry.uri ?? entry.name;
      if (!key || seen.has(key)) continue;
      seen.set(key, entry);
    }
  }
  return Array.from(seen.values());
}

function shortName(name: string): string {
  // Bazel encodes "test.outputs/outputs.zip" as "test.outputs__outputs.zip"
  // in the BES file name field. Strip the directory prefix for display.
  const idx = name.lastIndexOf("__");
  return idx >= 0 ? name.slice(idx + 2) : name;
}

function uriToDownloadHref(uri: string): string {
  // bytestream://cas.rbe:8980/blobs/<digest>/<size>
  //   -> route through the bb-portal /blob proxy that already serves
  //      test.log; reuse the same resolution path.
  // file:// or empty URIs are returned as-is; the browser refuses
  // file:// from a remote origin, which is the right behavior.
  if (uri.startsWith("bytestream://")) {
    return `/blob?uri=${encodeURIComponent(uri)}`;
  }
  return uri;
}
