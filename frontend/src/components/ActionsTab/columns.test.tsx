import type { ReactElement } from "react";
import { describe, expect, it } from "vitest";
import type { BazelInvocationActionFragment } from "@/graphql/__generated__/graphql";
import { getColumns } from "./columns";

describe("ActionsTab selectable filters", () => {
  const columns = getColumns(
    ["CppCompile", "GoLink"],
    ["k8-fastbuild", "k8-opt"],
  );

  it("filters by all selected action mnemonics", () => {
    const column = columns.find(({ key }) => key === "type");

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.filters).toEqual([
      { text: "CppCompile", value: "CppCompile" },
      { text: "GoLink", value: "GoLink" },
    ]);
    expect(column?.applyFilter?.(["CppCompile", "GoLink"])).toEqual([
      { typeIn: ["CppCompile", "GoLink"] },
    ]);
  });

  it("filters by all selected configurations", () => {
    const column = columns.find(({ key }) => key === "configurationMnemonic");

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.applyFilter?.(["k8-fastbuild", "k8-opt"])).toEqual([
      {
        hasConfigurationWith: [{ mnemonicIn: ["k8-fastbuild", "k8-opt"] }],
      },
    ]);
  });

  it("renders Action and primary output as separate links", () => {
    const action = {
      actionDigest: {
        rev2InstanceName: "projects/example",
        digestFunction: "SHA256",
        hash: "action-hash",
        sizeBytes: 145,
      },
      primaryOutput: "bazel-out/bin/output.txt",
      primaryOutputURI:
        "bytestream://cache.example.com/projects/example/blobs/output-hash/42",
    } as BazelInvocationActionFragment;

    const actionColumn = columns.find(({ key }) => key === "action");
    const actionLink = actionColumn?.render?.(
      undefined,
      action,
      0,
    ) as ReactElement<{ href: string }>;
    expect(actionLink.props.href).toBe(
      "/browser/projects/example/blobs/sha256/action/action-hash-145",
    );

    const primaryOutputColumn = columns.find(
      ({ key }) => key === "primaryOutput",
    );
    const primaryOutputLink = primaryOutputColumn?.render?.(
      undefined,
      action,
      0,
    ) as ReactElement<{
      title: string;
      children: ReactElement<{ href: string; children: string }>;
    }>;
    expect(primaryOutputLink.props.title).toBe("bazel-out/bin/output.txt");
    expect(primaryOutputLink.props.children.props.href).toBe(
      "/api/v1/servefile/projects/example/blobs/sha256/file/output-hash-42/output.txt",
    );
    expect(primaryOutputLink.props.children.props.children).toBe("output.txt");
  });
});
