import { Space, Typography } from "antd";
import Content from "@/components/Content";
import Uploader from "@/components/Uploader";
import { env } from "@/utils/env";
import styles from "./index.module.css";

const BuildInstructions: React.FC = () => {
  const bazelrcLines = `build --bes_backend=${env.grpcBackendUrl}\nbuild --bes_results_url=${window.location.origin}/bazel-invocations/`;
  return (
    <Space direction="vertical" size="large">
      <Typography.Text>
        Add the following lines to your{" "}
        <Typography.Text italic>.bazelrc</Typography.Text> to start sending
        build events to the service:
      </Typography.Text>
      <Space size="middle">
        <Typography.Text copyable={{ text: bazelrcLines }} />
        <pre style={{ textAlign: "left" }}>{bazelrcLines}</pre>
      </Space>
    </Space>
  );
};

const BepFileUploader: React.FC = () => {
  return (
    <Space direction="vertical" size="large">
      <Uploader
        label="Upload Build Event Protocol (BEP) files to analyze"
        description={
          <Typography.Text type="secondary">
            Upload one or more{" "}
            <Typography.Text type="secondary" italic>
              *.bep.ndjson
            </Typography.Text>{" "}
            file(s) produced with Bazel&apos;s{" "}
            <Typography.Text code>--build_event_json_file</Typography.Text> flag
            to analyze
          </Typography.Text>
        }
        action={"/api/v1/bep/upload"}
      />
    </Space>
  );
};

export function HomePage() {
  return (
    <Content
      content={
        <Space direction="vertical" size="large" className={styles.container}>
          {!!env.featureFlags?.home?.fileUpload && <BepFileUploader />}
          {!!env.featureFlags?.home?.instructions && <BuildInstructions />}
        </Space>
      }
    />
  );
}
