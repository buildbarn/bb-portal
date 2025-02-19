import { useGrpcClients } from "@/context/GrpcClientsContext";
import { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { BrowserPageParams } from "@/types/BrowserPageType";
import { useQuery } from "@tanstack/react-query";
import { Space, Spin, Typography } from "antd";
import BrowserPreviousExecutionsDisplay from "../BrowserPreviousExecutionsDisplay";
import PortalAlert from "../PortalAlert";

interface Params {
  browserPageParams: BrowserPageParams;
}

const BrowserPreviousExecutionsPage: React.FC<Params> = ({
  browserPageParams,
}) => {
  const { initialSizeClassCacheClient } = useGrpcClients();

  const reducedActionDigest = Digest.create(browserPageParams.digest);

  const { data, isPending, isError, error } = useQuery({
    queryKey: ["browserPreviousExecutionsPage", browserPageParams],
    queryFn: initialSizeClassCacheClient.getPreviousExecutionStats.bind(null, {
      digestFunction: browserPageParams.digestFunction,
      instanceName: browserPageParams.instanceName,
      reducedActionDigest: reducedActionDigest,
    }),
  });

  if (isError) {
    return (
      <PortalAlert
        className="error"
        message={
          <>
            <Typography.Text>
              There was a problem communicating with the backend server:
            </Typography.Text>
            <pre>{String(error)}</pre>
          </>
        }
      />
    );
  }

  if (isPending) {
    return <Spin />;
  }

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <Typography.Title level={2}>Previous executions stats</Typography.Title>
      <BrowserPreviousExecutionsDisplay
        browserParams={browserPageParams}
        previousExecutionStats={data}
        reducedActionDigest={reducedActionDigest}
        showTitle={false}
      />
    </Space>
  );
};

export default BrowserPreviousExecutionsPage;
