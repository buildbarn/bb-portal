import { expect, test, vi } from "vitest";
import { RequestMetadata } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { ProtobufTypeUrls } from "@/types/protobufTypeUrls";
import {
  anyToInvocationIdSearch,
  childInvocationNameSearch,
  type InvocationNameSearch,
  invocationIdSearchToAny,
  invocationIdToString,
  invocationNameFromSearch,
  rootInvocationNameSearch,
} from "./invocationName";

const sizeClassQueueName = {
  platformQueueName: {
    instanceNamePrefix: "foo/bar",
    platform: { properties: [{ name: "OSFamily", value: "linux" }] },
  },
  sizeClass: 1,
};

test("base64 round-trip", () => {
  const values = [
    new Uint8Array([]),
    new Uint8Array([0x00, 0x7f, 0x80, 0xc3, 0xff]),
    RequestMetadata.encode(
      RequestMetadata.fromPartial({ toolInvocationId: "abc" }),
    ).finish(),
  ];
  for (const value of values) {
    const search = anyToInvocationIdSearch({ typeUrl: "t", value });
    expect(search.typeUrl).toBe("t");
    expect(invocationIdSearchToAny(search)).toEqual({ typeUrl: "t", value });
  }
  expect(
    anyToInvocationIdSearch({ typeUrl: "t", value: new Uint8Array([]) }).value,
  ).toBe("");
  expect(
    anyToInvocationIdSearch({
      typeUrl: "t",
      value: new Uint8Array([0xfb, 0xff]),
    }).value,
  ).toBe("+/8=");
});

test("childInvocationNameSearch", () => {
  const root = rootInvocationNameSearch(sizeClassQueueName);
  expect(root).toEqual({ sizeClassQueueName, ids: [] });

  const child = childInvocationNameSearch(root, {
    typeUrl: "a",
    value: new Uint8Array([1]),
  });
  const grandchild = childInvocationNameSearch(child, {
    typeUrl: "b",
    value: new Uint8Array([2]),
  });
  expect(root.ids).toEqual([]);
  expect(child.ids).toEqual([{ typeUrl: "a", value: "AQ==" }]);
  expect(grandchild.ids).toEqual([
    { typeUrl: "a", value: "AQ==" },
    { typeUrl: "b", value: "Ag==" },
  ]);

  const parent: InvocationNameSearch = grandchild;
  expect(invocationNameFromSearch(parent)).toEqual({
    sizeClassQueueName,
    ids: [
      { typeUrl: "a", value: new Uint8Array([1]) },
      { typeUrl: "b", value: new Uint8Array([2]) },
    ],
  });
});

test("invocationIdToString", () => {
  expect(
    invocationIdToString({
      typeUrl: ProtobufTypeUrls.REQUEST_METADATA,
      value: RequestMetadata.encode(
        RequestMetadata.fromPartial({ toolInvocationId: "abc" }),
      ).finish(),
    }),
  ).toBe(
    JSON.stringify({
      "@type": ProtobufTypeUrls.REQUEST_METADATA,
      toolInvocationId: "abc",
    }),
  );
  expect(
    invocationIdToString({
      typeUrl: ProtobufTypeUrls.BACKGROUND_LEARNING,
      value: new Uint8Array([]),
    }),
  ).toBe(JSON.stringify({ "@type": ProtobufTypeUrls.BACKGROUND_LEARNING }));

  const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});
  expect(
    invocationIdToString({
      typeUrl: "type.googleapis.com/example.Unknown",
      value: new Uint8Array([1, 2, 3]),
    }),
  ).toBe(JSON.stringify({ "@type": "type.googleapis.com/example.Unknown" }));
  consoleError.mockRestore();
});
