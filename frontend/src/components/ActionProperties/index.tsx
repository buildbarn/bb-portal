import { Descriptions, Grid, theme } from "antd";
import CopyBbClientdActionButton from "@/components/BrowserActionGrid/CopyBbClientdActionButton";
import type { fetchBrowserActionGrid } from "@/components/BrowserActionGrid/fetch";
import type { BrowserPageParams } from "@/types/BrowserPageParams";
import { readableDurationFromProtobufDuration } from "@/utils/time";
import { styleMap } from "../BrowserDirectory/utils";
import PropertyTagList from "../PropertyTagList";
import ComparePropertyTagList from "../PropertyTagList/comparePropertyTagList";

const { useToken } = theme;

interface Params {
  browserPageParams: BrowserPageParams;
  actionData: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  compareActionData?: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  mergedMode?: boolean;
}
const InnerActionProperties: React.FC<Params> = ({
  browserPageParams,
  actionData,
  compareActionData,
  mergedMode, // For comparing actions
}) => {
  const screens = Grid.useBreakpoint();
  const { token } = useToken();
  return (
    <>
      <Descriptions
        layout={screens.md ? "horizontal" : "vertical"}
        column={1}
        size="small"
        bordered
        styles={{ label: { width: "25%" }, content: { width: "75%" } }}
      >
        {actionData.action.timeout && (
          <Descriptions.Item
            label="Timeout:"
            style={
              compareActionData &&
              actionData.action.timeout !== compareActionData?.action.timeout
                ? styleMap("diff_with_borders", token)
                : {}
            }
          >
            {readableDurationFromProtobufDuration(actionData.action.timeout)}
          </Descriptions.Item>
        )}
        <Descriptions.Item
          label="Do not cache"
          style={
            compareActionData &&
            actionData.action.doNotCache !==
              compareActionData?.action.doNotCache
              ? styleMap("diff_with_borders", token)
              : {}
          }
        >
          {actionData.action.doNotCache ? "Yes" : "No"}
        </Descriptions.Item>
        {actionData.action.platform &&
          (compareActionData?.action.platform ? (
            <Descriptions.Item label="Platform properties">
              <ComparePropertyTagList
                propertyList={actionData.action.platform.properties}
                comparePropertyList={
                  compareActionData?.action.platform?.properties
                }
                mergedMode={mergedMode}
              />
            </Descriptions.Item>
          ) : (
            <Descriptions.Item label="Platform properties">
              <PropertyTagList
                propertyList={actionData.action.platform.properties}
              />
            </Descriptions.Item>
          ))}
      </Descriptions>
      {actionData.action.commandDigest &&
        actionData.action.inputRootDigest &&
        !mergedMode && (
          <CopyBbClientdActionButton
            instanceName={browserPageParams.instanceName}
            digestFunction={browserPageParams.digestFunction}
            actionDigest={actionData.actionDigest}
            commandDigest={actionData.action.commandDigest}
            inputRootDigest={actionData.action.inputRootDigest}
          />
        )}
    </>
  );
};

interface ActionPropertiesParams {
  browserPageParams: BrowserPageParams;
  actionData: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
}
const ActionProperties = ({
  browserPageParams,
  actionData,
}: ActionPropertiesParams) => {
  return (
    <InnerActionProperties
      browserPageParams={browserPageParams}
      actionData={actionData}
    />
  );
};
interface ActionPropertiesDeltaParams {
  browserPageParams: BrowserPageParams;
  actionData: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  compareActionData: Awaited<ReturnType<typeof fetchBrowserActionGrid>>;
  mergedMode: boolean;
}
const ActionPropertiesDelta = ({
  browserPageParams,
  actionData,
  compareActionData,
  mergedMode,
}: ActionPropertiesDeltaParams) => {
  return (
    <InnerActionProperties
      browserPageParams={browserPageParams}
      actionData={actionData}
      compareActionData={compareActionData}
      mergedMode={mergedMode}
    />
  );
};

export { ActionProperties, ActionPropertiesDelta };
