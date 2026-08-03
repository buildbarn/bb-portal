import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ClientError, Status } from "nice-grpc-web";
import { act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, describe, expect, it, vi } from "vitest";
import { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import BrowserDirectory from ".";

(
  globalThis as typeof globalThis & { IS_REACT_ACT_ENVIRONMENT: boolean }
).IS_REACT_ACT_ENVIRONMENT = true;

const { read } = vi.hoisted(() => ({ read: vi.fn() }));

vi.mock("@/grpc/casByteStreamClient", () => ({
  casByteStreamClient: { read },
}));

afterEach(() => {
  read.mockReset();
});

describe("BrowserDirectory", () => {
  it("explains a missing action input tree and hides unusable controls", async () => {
    read.mockImplementation(() => {
      throw new ClientError(
        "/google.bytestream.ByteStream/Read",
        Status.UNKNOWN,
        "Shard 0: Object not found",
      );
    });

    const container = document.createElement("div");
    const root = createRoot(container);
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    await act(async () => {
      root.render(
        <QueryClientProvider client={queryClient}>
          <BrowserDirectory
            instanceName=""
            digestFunction={DigestFunction_Value.SHA256}
            inputRootDigest={{ hash: "missing-input-root", sizeBytes: "83" }}
            fileSystemAccessProfile={undefined}
            fileSystemAccessProfileReference={undefined}
            notFoundMessage="Input tree not uploaded"
            notFoundDescription="The input root is not available in the CAS."
          />
        </QueryClientProvider>,
      );
    });

    await vi.waitFor(() => {
      expect(container.textContent).toContain("Input tree not uploaded");
    });

    expect(container.textContent).toContain(
      "The input root is not available in the CAS.",
    );
    expect(container.textContent).not.toContain("Download as tarball");
    expect(container.textContent).not.toContain(
      "Copy bb_clientd path to clipboard",
    );

    await act(async () => root.unmount());
  });
});
