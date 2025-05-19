"use client";

import Content from "@/components/Content";
import PageDisabled from "@/components/PageDisabled";
import PortalCard from "@/components/PortalCard";
import { FeatureType, isFeatureEnabled } from "@/utils/isFeatureEnabled";
import { LayoutOutlined } from "@ant-design/icons";
import { Space, Typography } from "antd";
import { notFound } from "next/navigation";
import type React from "react";
import styles from "./page.module.css";

const Page: React.FC = () => {
  if (!isFeatureEnabled(FeatureType.BROWSER)) {
    return <PageDisabled />;
  }

  return (
    <Content
      content={
        <Space direction="vertical" size="middle" style={{ display: "flex" }}>
          <PortalCard
            icon={<LayoutOutlined />}
            titleBits={[<span key="title">Browser</span>]}
          >
            <Typography.Title level={2}>Welcome</Typography.Title>
            <Typography.Text>
              This page allows you to display objects stored in the Content
              Addressable Storage (CAS) and Action Cache (AC) as defined by the{" "}
              <a href="https://github.com/bazelbuild/remote-apis">
                Remote Execution API
              </a>
              . Objects in these data stores have hard to guess identifiers and
              the Remote Execution API provides no functions for iterating over
              them. One may therefore only access this service in a meaningful
              way by visiting automatically generated URLs pointing to this
              page. Tools that are part of Buildbarn will generate these URLs
              where applicable.
            </Typography.Text>

            <Typography.Paragraph>
              <Typography.Text>
                This service supports the following URL schemes:
              </Typography.Text>

              <ul className={styles.welcomeList}>
                <li>
                  <p>
                    <pre>
                      {
                        "/browser/${instance_name}/blobs/${digest_function}/action/${hash}-${size_bytes}/"
                      }
                    </pre>
                    Displays information about an Action and its associated
                    Command stored in the CAS. If available, displays
                    information about the Action&apos;s associated ActionResult
                    stored in the AC.
                  </p>
                </li>
                <li>
                  <p>
                    <pre>
                      {
                        "/browser/${instance_name}/blobs/${digest_function}/command/${hash}-${size_bytes}/"
                      }
                    </pre>
                    Displays information about a Command stored in the CAS.
                  </p>
                </li>
                <li>
                  <p>
                    <pre>
                      {
                        "/browser/${instance_name}/blobs/${digest_function}/directory/${hash}-${size_bytes}/"
                      }
                    </pre>
                    Displays information about a Directory stored in the CAS.
                  </p>
                </li>
                <li>
                  <p>
                    <pre>
                      {
                        "/api/v1/servefile/${instance_name}/blobs/${digest_function}/file/${hash}-${size_bytes}/${filename}"
                      }
                    </pre>
                    Serves a file stored in the CAS.
                  </p>
                </li>
                <li>
                  <p>
                    <pre>
                      {
                        "/browser/${instance_name}/blobs/${digest_function}/historical_execute_response/${hash}-${size_bytes}/"
                      }
                    </pre>
                    Extension: displays information about an ActionResult that
                    was not permitted to be stored in the AC, but was stored in
                    the CAS instead. Buildbarn stores ActionResult messages for
                    failed build actions in the CAS.
                  </p>
                </li>
                <li>
                  <p>
                    <pre>
                      {
                        "/browser/${instance_name}/blobs/${digest_function}/previous_execution_stats/${hash}-${size_bytes}/"
                      }
                    </pre>
                    Extension: displays information about outcomes of previous
                    executions of similar actions. This information is extracted
                    from Buildbarn&apos;s Initial Size Class Cache (ISCC).
                  </p>
                </li>
                <li>
                  <p>
                    <pre>
                      {
                        "/browser/${instance_name}/blobs/${digest_function}/tree/${hash}-${size_bytes}/${subdirectory}/"
                      }
                    </pre>
                    Displays information about a Directory contained in a Tree
                    stored in the CAS.
                  </p>
                </li>
              </ul>
            </Typography.Paragraph>
          </PortalCard>
        </Space>
      }
    />
  );
};

export default Page;
