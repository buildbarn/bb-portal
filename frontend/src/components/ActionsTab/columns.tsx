import { SearchOutlined } from "@ant-design/icons";
import { Tag, Tooltip, Typography } from "antd";
import type { FilterDropdownProps, FilterValue } from "antd/es/table/interface";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type {
  ActionWhereInput,
  BazelInvocationActionFragment,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";
import {
  generateActionUrlFromGraphqlDigest,
  generateFileUrlFromBepURI,
} from "@/utils/urlGenerator";
import styles from "../../theme/theme.module.css";
import PortalDuration from "../PortalDuration";
import { getActionCacheStatus } from "./cache";
import { getActionExecutionKind } from "./execution";

const searchFilter = (
  placeholder: string,
  toWhere: (value: string) => ActionWhereInput,
) => ({
  filterSearch: true,
  filterDropdown: (filterProps: FilterDropdownProps) => (
    <SearchWidget {...filterProps} placeholder={placeholder} />
  ),
  filterIcon: (filtered: boolean) => (
    <SearchFilterIcon icon={<SearchOutlined />} filtered={filtered} />
  ),
  applyFilter: (value: FilterValue) => {
    if (!value || value.length === 0) {
      return undefined;
    }
    return [toWhere(value[0] as string)];
  },
});

const selectableFilter = <Value extends string>(
  options: readonly Value[],
  toWhere: (values: Value[]) => ActionWhereInput,
) => {
  const uniqueOptions = Array.from(new Set(options));
  return {
    filters: uniqueOptions.map((value) => ({ text: value, value })),
    filterMultiple: true,
    filterSearch: true,
    applyFilter: (value: FilterValue) => {
      const values = value.filter(
        (selectedValue): selectedValue is Value =>
          typeof selectedValue === "string" &&
          uniqueOptions.includes(selectedValue as Value),
      );
      return values.length > 0 ? [toWhere(values)] : undefined;
    },
  };
};

const configurationTooltip = (
  configuration: NonNullable<BazelInvocationActionFragment["configuration"]>,
) => {
  const makeVariables =
    configuration.makeVariables &&
    typeof configuration.makeVariables === "object" &&
    !Array.isArray(configuration.makeVariables)
      ? Object.entries(configuration.makeVariables)
      : [];

  return (
    <>
      <b>Configuration ID:</b> <code>{configuration.configurationID}</code>
      <span style={{ display: "block" }}>
        <b>Mnemonic:</b> <code>{configuration.mnemonic || "-"}</code>
      </span>
      <span style={{ display: "block" }}>
        <b>Platform:</b> <code>{configuration.platformName || "-"}</code>
      </span>
      <span style={{ display: "block" }}>
        <b>CPU:</b> <code>{configuration.cpu || "-"}</code>
      </span>
      {makeVariables.length > 0 && (
        <>
          <span style={{ display: "block" }}>
            <b>Make variables:</b>
          </span>
          {makeVariables.map(([key, value]) => (
            <span key={key} style={{ display: "block", paddingLeft: "1em" }}>
              <code>
                {key}=
                {typeof value === "string" ? value : JSON.stringify(value)}
              </code>
            </span>
          ))}
        </>
      )}
    </>
  );
};

const actionStatuses = ["Succeeded", "Failed"] as const;

const actionStatusPredicates: Record<
  (typeof actionStatuses)[number],
  ActionWhereInput
> = {
  Succeeded: { success: true },
  Failed: { or: [{ success: false }, { successIsNil: true }] },
};

const executionKinds = ["Remote", "Local", "Internal", "Unknown"] as const;

const executionPredicates: Record<
  (typeof executionKinds)[number],
  ActionWhereInput
> = {
  Remote: { runnerIn: ["remote", "remote cache hit"] },
  Local: {
    runnerNotNil: true,
    runnerNEQ: "",
    runnerNotIn: ["remote", "remote cache hit", "internal"],
  },
  Internal: { runner: "internal" },
  Unknown: { or: [{ runnerIsNil: true }, { runner: "" }] },
};

const cacheResults = [
  "Remote hit",
  "Disk hit",
  "Hit",
  "No hit",
  "Unknown",
] as const;

const cacheResultPredicates: Record<
  (typeof cacheResults)[number],
  ActionWhereInput
> = {
  "Remote hit": { cacheHit: true, runner: "remote cache hit" },
  "Disk hit": { cacheHit: true, runner: "disk cache hit" },
  Hit: {
    cacheHit: true,
    or: [
      { runnerIsNil: true },
      { runnerNotIn: ["remote cache hit", "disk cache hit"] },
    ],
  },
  "No hit": { cacheHit: false },
  Unknown: { cacheHitIsNil: true },
};

export const getColumns = (
  actionMnemonics: string[],
  configurationMnemonics: string[],
): TableColumnTypeWithFilter<
  BazelInvocationActionFragment,
  ActionWhereInput
>[] => [
  {
    key: "success",
    title: "Status",
    width: 100,
    render: (_, action) =>
      action.success ? (
        <Tag color="success">Succeeded</Tag>
      ) : (
        <Tag color="error">Failed</Tag>
      ),
    ...selectableFilter(actionStatuses, (values) => ({
      or: values.map((value) => actionStatusPredicates[value]),
    })),
  },
  {
    key: "duration",
    title: "Duration",
    width: 120,
    render: (_, action) => (
      <span className={styles.numberFormat}>
        {action.startTime && action.endTime ? (
          <PortalDuration
            from={action.startTime}
            to={action.endTime}
            formatConfig={{ smallestUnit: "ms" }}
          />
        ) : (
          "No spawn timing"
        )}
      </span>
    ),
  },
  {
    key: "execution",
    title: "Execution",
    width: 110,
    render: (_, action) => {
      const executionKind = getActionExecutionKind(action.runner);
      const colors = {
        Remote: "processing",
        Local: "warning",
        Internal: "purple",
        Unknown: "default",
      } as const;
      return (
        <Tooltip title={action.runner || "Not reported by Bazel"}>
          <Tag color={colors[executionKind]}>{executionKind}</Tag>
        </Tooltip>
      );
    },
    ...selectableFilter(executionKinds, (values) => ({
      or: values.map((value) => executionPredicates[value]),
    })),
  },
  {
    key: "cacheHit",
    title: "Cache result",
    width: 110,
    render: (_, action) => {
      const cacheStatus = getActionCacheStatus(action.cacheHit, action.runner);
      return (
        <Tooltip title={cacheStatus.description}>
          <Tag color={cacheStatus.color}>{cacheStatus.label}</Tag>
        </Tooltip>
      );
    },
    ...selectableFilter(cacheResults, (values) => ({
      or: values.map((value) => cacheResultPredicates[value]),
    })),
  },
  {
    key: "type",
    title: "Action mnemonic",
    dataIndex: "type",
    width: 180,
    ellipsis: true,
    ...selectableFilter(actionMnemonics, (values) => ({
      typeIn: values,
    })),
  },
  {
    key: "configurationMnemonic",
    title: "Configuration",
    width: 180,
    ellipsis: true,
    render: (_, action) => {
      if (!action.configuration) {
        return undefined;
      }
      return (
        <Tooltip
          placement="topLeft"
          styles={{ root: { maxWidth: "50vw" } }}
          title={configurationTooltip(action.configuration)}
        >
          <span>
            {action.configuration.mnemonic ||
              action.configuration.configurationID}
          </span>
        </Tooltip>
      );
    },
    ...selectableFilter(configurationMnemonics, (values) => ({
      hasConfigurationWith: [{ mnemonicIn: values }],
    })),
  },
  {
    key: "label",
    title: "Label",
    dataIndex: "label",
    ellipsis: true,
    ...searchFilter("Action label...", (value) => ({
      labelContainsFold: value,
    })),
  },
  {
    key: "action",
    title: "Action cache",
    width: 120,
    render: (_, action) => {
      const href = generateActionUrlFromGraphqlDigest(action.actionDigest);
      return href ? (
        <Typography.Link href={href}>View action</Typography.Link>
      ) : undefined;
    },
  },
  {
    key: "primaryOutput",
    title: "Primary output",
    ellipsis: true,
    render: (_, action) => {
      if (!action.primaryOutput) {
        return undefined;
      }
      const fileName =
        action.primaryOutput.split("/").pop() || action.primaryOutput;
      const href = generateFileUrlFromBepURI(
        action.primaryOutputURI,
        action.primaryOutput,
      );
      return (
        <Tooltip title={action.primaryOutput}>
          {href ? (
            <Typography.Link href={href}>{fileName}</Typography.Link>
          ) : (
            fileName
          )}
        </Tooltip>
      );
    },
  },
];
