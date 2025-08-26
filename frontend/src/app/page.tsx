"use client";
import Content from "@/components/Content";
import Uploader from "@/components/Uploader";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { Divider, Space, Typography } from "antd";
import { env } from "next-runtime-env";
import styles from "./page.module.css";

export default function Home() {
  const bazelrcLines = `build --bes_backend=${env("NEXT_PUBLIC_BES_GRPC_BACKEND_URL")}\nbuild --bes_results_url=${window.location.origin}/bazel-invocations/`;

  return (
    <Content
      content={
        <div className={styles.container}>
          <Space direction="vertical" size="large">
            <Typography.Title level={1} className={styles.item}>
              Welcome to the {env("NEXT_PUBLIC_COMPANY_NAME")} Buildbarn Portal
            </Typography.Title>
            <Typography.Title level={5} className={styles.item}>
              Providing insights into Bazel build outputs
            </Typography.Title>
            {isFeatureEnabled(FeatureType.BES) && (
              <>
                {isFeatureEnabled(FeatureType.BEP_UPLOAD) && (
                  <>
                    <Divider />
                    <Uploader
                      label="Upload Build Event Protocol (BEP) files to analyze"
                      description={
                        <Typography.Text type="secondary">
                          Upload one or more{" "}
                          <Typography.Text type="secondary" italic>
                            *.bep.ndjson
                          </Typography.Text>{" "}
                          file(s) produced with Bazel&apos;s{" "}
                          <Typography.Text code>
                            --build_event_json_file
                          </Typography.Text>{" "}
                          flag to analyze
                        </Typography.Text>
                      }
                      action={`${env("NEXT_PUBLIC_BES_BACKEND_URL")}/api/v1/bep/upload`}
                    />
                    <Divider />
                  </>
                )}
                <Typography.Text>
                  Add the following lines to your{" "}
                  <Typography.Text italic>.bazelrc</Typography.Text> to start
                  sending build events to the service:
                </Typography.Text>
                <Space size="middle">
                  <Typography.Text copyable={{ text: bazelrcLines }} />
                  <pre style={{ textAlign: "left" }}>{bazelrcLines}</pre>
                </Space>
              </>
            )}
          </Space>
        </div>
      }
    />
  );
}
