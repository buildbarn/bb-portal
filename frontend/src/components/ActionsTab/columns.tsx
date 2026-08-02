import { SearchOutlined } from "@ant-design/icons";
import { Tag, Typography } from "antd";
import type { FilterDropdownProps, FilterValue } from "antd/es/table/interface";
import { SearchFilterIcon, SearchWidget } from "@/components/SearchWidgets";
import type {
  ActionWhereInput,
  BazelInvocationActionFragment,
} from "@/graphql/__generated__/graphql";
import type { TableColumnTypeWithFilter } from "@/types/TableColumnTypeWithFilter";
import { generateFileUrlFromBepURI } from "@/utils/urlGenerator";
import styles from "../../theme/theme.module.css";
import PortalDuration from "../PortalDuration";

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

export const getColumns = (): TableColumnTypeWithFilter<
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
    key: "type",
    title: "Action mnemonic",
    dataIndex: "type",
    width: 180,
    ellipsis: true,
    ...searchFilter("Action mnemonic...", (value) => ({
      typeContainsFold: value,
    })),
  },
  {
    key: "configurationMnemonic",
    title: "Configuration",
    width: 180,
    ellipsis: true,
    render: (_, action) => action.configuration?.mnemonic,
    ...searchFilter("Configuration mnemonic...", (value) => ({
      hasConfigurationWith: [{ mnemonicContainsFold: value }],
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
      const href = generateFileUrlFromBepURI(
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
