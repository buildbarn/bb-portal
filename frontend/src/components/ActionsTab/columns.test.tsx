import type { ReactElement } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import type { BazelInvocationActionExecutionFragment } from "@/graphql/__generated__/graphql";
import { ActionExecutionWhereInputSchema } from "@/graphql/__generated__/zod";
import { getColumns } from "./columns";

describe("ActionsTab selectable filters", () => {
  const columns = getColumns(
    ["CppCompile", "GoLink"],
    ["k8-fastbuild", "k8-opt"],
  );

  it("filters by all selected action statuses", () => {
    const column = columns.find(({ key }) => key === "success");

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.filters).toEqual([
      { text: "Succeeded", value: "Succeeded" },
      { text: "Failed", value: "Failed" },
    ]);
    expect(column?.applyFilter?.(["Succeeded", "Failed"])).toEqual([
      {
        or: [
          { success: true },
          { or: [{ success: false }, { successIsNil: true }] },
        ],
      },
    ]);
  });

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

  it("removes duplicate configuration filter options", () => {
    const configurationColumn = getColumns(
      [],
      ["k8-fastbuild", "k8-opt", "k8-fastbuild"],
    ).find(({ key }) => key === "configurationMnemonic");

    expect(configurationColumn?.filters).toEqual([
      { text: "k8-fastbuild", value: "k8-fastbuild" },
      { text: "k8-opt", value: "k8-opt" },
    ]);
  });

  it("shows configuration details on hover", () => {
    const configurationColumn = columns.find(
      ({ key }) => key === "configurationMnemonic",
    );
    const configurationCell = configurationColumn?.render?.(
      undefined,
      {
        configuration: {
          configurationID: "configuration-id",
          mnemonic: "k8-fastbuild",
          platformName: "linux-arm64",
          cpu: "arm64",
          makeVariables: { COMPILATION_MODE: "fastbuild" },
        },
      } as BazelInvocationActionExecutionFragment,
      0,
    ) as ReactElement<{
      title: React.ReactNode;
      children: ReactElement<{ children: string }>;
    }>;

    expect(configurationCell.props.children.props.children).toBe(
      "k8-fastbuild",
    );
    const tooltip = renderToStaticMarkup(configurationCell.props.title);
    expect(tooltip).toContain("Configuration ID:");
    expect(tooltip).toContain("configuration-id");
    expect(tooltip).toContain("linux-arm64");
    expect(tooltip).toContain("arm64");
    expect(tooltip).toContain("COMPILATION_MODE=fastbuild");
  });

  it("filters by all selected execution kinds", () => {
    const column = columns.find(({ key }) => key === "execution");
    const filter = column?.applyFilter?.(["Local", "Unknown"]);

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.filters).toEqual([
      { text: "Remote", value: "Remote" },
      { text: "Local", value: "Local" },
      { text: "Internal", value: "Internal" },
      { text: "Unknown", value: "Unknown" },
    ]);
    expect(filter).toEqual([
      {
        or: [
          {
            runnerNotNil: true,
            runnerNEQ: "",
            runnerNotIn: ["remote", "remote cache hit", "internal"],
          },
          { or: [{ runnerIsNil: true }, { runner: "" }] },
        ],
      },
    ]);
    expect(
      ActionExecutionWhereInputSchema().partial().safeParse(filter?.[0])
        .success,
    ).toBe(true);
  });

  it("filters by all selected cache results", () => {
    const column = columns.find(({ key }) => key === "cacheHit");

    expect(column?.filterSearch).toBe(true);
    expect(column?.filterMultiple).toBe(true);
    expect(column?.filters).toEqual([
      { text: "Remote hit", value: "Remote hit" },
      { text: "Disk hit", value: "Disk hit" },
      { text: "Hit", value: "Hit" },
      { text: "No hit", value: "No hit" },
      { text: "Unknown", value: "Unknown" },
    ]);
    expect(column?.applyFilter?.(["Remote hit", "Hit", "Unknown"])).toEqual([
      {
        or: [
          { cacheHit: true, runner: "remote cache hit" },
          {
            cacheHit: true,
            or: [
              { runnerIsNil: true },
              {
                runnerNotIn: ["remote cache hit", "disk cache hit"],
              },
            ],
          },
          { cacheHitIsNil: true },
        ],
      },
    ]);
  });

  it("renders cache-hit metadata in its own column", () => {
    const cacheColumn = columns.find(({ key }) => key === "cacheHit");
    const cacheCell = cacheColumn?.render?.(
      undefined,
      {
        cacheHit: true,
        runner: "remote cache hit",
      } as BazelInvocationActionExecutionFragment,
      0,
    ) as ReactElement<{
      title: string;
      children: ReactElement<{ color: string; children: string }>;
    }>;

    expect(cacheCell.props.title).toBe("Bazel reported a remote cache hit");
    expect(cacheCell.props.children.props.color).toBe("success");
    expect(cacheCell.props.children.props.children).toBe("Remote hit");
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
    } as BazelInvocationActionExecutionFragment;

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
