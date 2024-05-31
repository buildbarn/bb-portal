import React from "react";
import {
  Button,
  Input,
  Space,
  Divider,
  DatePicker,
  TimeRangePickerProps,
} from "antd";
import { FilterDropdownProps } from "antd/es/table/interface";
import dayjs from "dayjs";
import { blue } from "@ant-design/colors";
import styles from "@/components/SearchWidgets/index.module.css";

interface SearchFilterIconProps {
  icon: React.ReactNode;
  filtered: boolean;
}

export const SearchFilterIcon: React.FC<SearchFilterIconProps> = ({
  icon,
  filtered,
}) => (
  <span style={{ color: filtered ? blue.primary : undefined }}>{icon}</span>
);

const ResetButton: React.FC<{ clearFilters: () => void }> = ({
  clearFilters,
}) => (
  <Button
    onClick={() => {
      clearFilters();
    }}
    size="small"
    type="link"
  >
    Reset
  </Button>
);

interface TimeRangeSelectorProps extends FilterDropdownProps {
  placeholder?: string;
  options?: readonly string[];
  renderOption?: (option: string) => React.ReactNode;
}

const timeRangePresets: TimeRangePickerProps["presets"] = [
  { label: "Starting Now", value: [dayjs(), null] },
  { label: "Until Now", value: [null, dayjs()] },
  { label: "Since Beginning of Day", value: [dayjs().startOf("day"), null] },
  { label: "Since Beginning of Week", value: [dayjs().startOf("week"), null] },
  {
    label: "Since Beginning of Month",
    value: [dayjs().startOf("month"), null],
  },
  { label: "Since Beginning of Year", value: [dayjs().startOf("year"), null] },
  { label: "Since 1 Hour Ago", value: [dayjs().add(-1, "hour"), null] },
  { label: "Since 1 Day Ago", value: [dayjs().add(-1, "day"), null] },
  { label: "Since 1 Week Ago", value: [dayjs().add(-1, "week"), null] },
  { label: "Since 1 Month Ago", value: [dayjs().add(-1, "month"), null] },
  { label: "Since 1 Year Ago", value: [dayjs().add(-1, "year"), null] },
  { label: "Past Hour", value: [dayjs().add(-1, "hour"), dayjs()] },
  { label: "Past Day", value: [dayjs().add(-1, "day"), dayjs()] },
  { label: "Past Week", value: [dayjs().add(-1, "week"), dayjs()] },
  { label: "Past Month", value: [dayjs().add(-1, "month"), dayjs()] },
  { label: "Past Year", value: [dayjs().add(-1, "year"), dayjs()] },
];

export const TimeRangeSelector: React.FC<TimeRangeSelectorProps> = ({
  setSelectedKeys,
  confirm,
}) => {
  return (
    <DatePicker.RangePicker
      allowEmpty={[true, true]}
      presets={timeRangePresets}
      showTime
      format="YYYY-MM-DD hh:mm:ss A"
      onChange={(dates) => {
        if (dates?.length === 2) {
          const from = dates[0] ? dates[0].toISOString() : "";
          const to = dates[1] ? dates[1].toISOString() : "";
          setSelectedKeys([from, to]);
        } else {
          setSelectedKeys([]);
        }
        confirm();
      }}
      className={[
        styles.searchWidgetInput,
        styles.searchWidgetRangePickerInput,
      ].join(" ")}
    />
  );
};

interface InputFilterProps extends FilterDropdownProps {
  placeholder: string;
  testID?: string;
}

export const SearchWidget: React.FC<InputFilterProps> = ({
  placeholder,
  selectedKeys,
  setSelectedKeys,
  clearFilters,
  confirm,
  testID,
}) => {
  if (!clearFilters) {
    // Pretty sure ant-design's types are just too loose, don't expect to ever be called without this callback
    return <p>Selections unavailable... missing callback</p>;
  }
  return (
    <div data-testid={testID}>
      <Space direction="vertical">
        <Input
          placeholder={placeholder}
          value={selectedKeys && selectedKeys[0]}
          onChange={(e) =>
            setSelectedKeys(e.target.value ? [e.target.value] : [])
          }
          onBlur={() => confirm()}
          onPressEnter={() => confirm()}
          className={[
            styles.searchWidgetInput,
            styles.searchWidgetTextInput,
          ].join(" ")}
        />
        <Divider className={styles.searchWidgetDivider} />
        <div className={styles.searchWidgetButtons}>
          <ResetButton clearFilters={clearFilters} />
          <div className={styles.searchWidgetButtonsSpacing} />
          <Button type="primary" size="small" onClick={() => confirm()}>
            OK
          </Button>
        </div>
      </Space>
    </div>
  );
};

export default SearchWidget;
