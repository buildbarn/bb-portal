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

const selectableFilter = (
  options: string[],
  toWhere: (values: string[]) => ActionWhereInput,
) => ({
  filters: options.map((value) => ({ text: value, value })),
  filterMultiple: true,
  filterSearch: true,
  applyFilter: (value: FilterValue) => {
    const values = value.filter(
      (selectedValue): selectedValue is string =>
        typeof selectedValue === "string",
    );
    return values.length > 0 ? [toWhere(values)] : undefined;
  },
});

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
    render: (_, action) => action.configuration?.mnemonic,
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
    key: "primaryOutput",
    title: "Primary output",
    ellipsis: true,
    render: (_, action) => {
      if (!action.primaryOutput) {
        return undefined;
      }
      const href =
        generateActionUrlFromGraphqlDigest(action.actionDigest) ??
        generateFileUrlFromBepURI(
          action.primaryOutputURI,
          action.primaryOutput,
        );
      return href ? (
        <Typography.Link href={href}>{action.primaryOutput}</Typography.Link>
      ) : (
        action.primaryOutput
      );
    },
  },
];
