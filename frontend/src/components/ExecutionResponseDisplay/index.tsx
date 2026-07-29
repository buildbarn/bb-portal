import { CodeFilled } from "@ant-design/icons";
import { ExecuteResponse } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { protobufToObjectWithTypeField } from "@/utils/protobufToObject";
import PortalCard from "../PortalCard";

interface Props {
  executeResponse: ExecuteResponse;
}

const ExecuteResponseDisplay: React.FC<Props> = ({ executeResponse }) => {
  const auxiliaryMetadata =
    executeResponse?.result?.executionMetadata?.auxiliaryMetadata.map(
      (value) => {
        return protobufToObjectWithTypeField(value, false);
      },
    );

  return (
    <PortalCard titleBits={["Execute Response"]} icon={<CodeFilled />}>
      <pre>
        {
          // `ts-proto` currently does not support JSON string
          // encoding of well-known type google.protobuf.Duration
          JSON.stringify(
            ExecuteResponse.toJSON(executeResponse),
            (key, val) => {
              if (key === "auxiliaryMetadata") {
                return auxiliaryMetadata;
              }
              return val;
            },
            1,
          )
        }
      </pre>
    </PortalCard>
  );
};

export default ExecuteResponseDisplay;
