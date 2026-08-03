import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it, vi } from "vitest";
import { MessageContext } from "@/context/MessageContext";
import {
  DigestFunction_Value,
  ExecuteResponse,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { BrowserPageType } from "@/types/BrowserPageType";
import BrowserResultDescription from ".";

const browserPageParams = {
  instanceName: "",
  digestFunction: DigestFunction_Value.SHA256,
  browserPageType: BrowserPageType.Action,
  digest: { hash: "action", sizeBytes: "10" },
  otherParams: [],
};

const renderResult = (executeResponse: ExecuteResponse) =>
  renderToStaticMarkup(
    <MessageContext.Provider
      value={{
        messageApi: {} as never,
        copyToClipboard: vi.fn(),
      }}
    >
      <BrowserResultDescription
        browserPageParams={browserPageParams}
        executeResponse={executeResponse}
        posixResourceUsage={undefined}
      />
    </MessageContext.Provider>,
  );

describe("BrowserResultDescription console output", () => {
  it("shows explicit empty output states without download buttons", () => {
    const emptyDigest = {
      hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      sizeBytes: "0",
    };
    const html = renderResult(
      ExecuteResponse.fromPartial({
        result: {
          exitCode: 0,
          stdoutDigest: emptyDigest,
          stderrDigest: emptyDigest,
        },
      }),
    );

    expect(html).toContain("The action produced no standard output.");
    expect(html).toContain("The action produced no standard error.");
    expect(html).not.toContain("Download Log");
  });

  it("shows when no console logs were uploaded", () => {
    const html = renderResult(
      ExecuteResponse.fromPartial({ result: { exitCode: 0 } }),
    );

    expect(html).toContain("No standard output log was uploaded.");
    expect(html).toContain("No standard error log was uploaded.");
    expect(html).not.toContain("Download Log");
  });
});
