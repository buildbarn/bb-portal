import type { BazelInvocationInfoFragment } from "@/graphql/__generated__/graphql";
import { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { generateFileUrl } from "@/utils/urlGenerator";
import {
  DownOutlined,
  DownloadOutlined,
  ProjectOutlined,
} from "@ant-design/icons";
import { Button, Dropdown, type MenuProps, Space } from "antd";

type Profile = NonNullable<BazelInvocationInfoFragment["profile"]>;

const PERFETTO_URL = "https://ui.perfetto.dev";

const getProfileUrl = (
  instanceName: string | undefined,
  profile: Profile,
): string => {
  return generateFileUrl(
    instanceName,
    DigestFunction_Value.SHA256,
    {
      hash: profile.digest,
      sizeBytes: profile.sizeInBytes.toString(),
    },
    profile.name,
  );
};

const fetchProfileFile = async (
  instanceName: string | undefined,
  profile: Profile,
): Promise<ArrayBuffer> => {
  const res = await fetch(getProfileUrl(instanceName, profile));
  if (!res.ok) {
    return Promise.reject(`Failed to download profile file: ${res.statusText}`);
  }
  return res.arrayBuffer();
};

const waitForPerfettoToLoad = async (handle: Window) => {
  const timer = setInterval(
    () => handle.postMessage("PING", PERFETTO_URL),
    100,
  );

  // Wait for the Perfetto UI to respond with 'PONG'
  await new Promise<void>((resolve) => {
    const listener = (evt: MessageEvent) => {
      if (evt.data !== "PONG") return;
      window.removeEventListener("message", listener);
      resolve();
    };
    window.addEventListener("message", listener);
  });

  window.clearInterval(timer);
};

function openPerfetto(
  instanceName: string | undefined,
  profile: Profile,
  invocationID: string,
) {
  const handle = window.open(PERFETTO_URL);
  if (!handle) {
    console.error("Failed to open new window for Perfetto UI");
    return;
  }

  Promise.all([
    fetchProfileFile(instanceName, profile),
    waitForPerfettoToLoad(handle),
  ])
    .then((values) => {
      handle.postMessage(
        {
          perfetto: {
            buffer: values[0],
            title: invocationID,
            fileName: profile.name,
          },
        },
        PERFETTO_URL,
      );
    })
    .catch((error) => {
      console.error("Error opening Perfetto:", error);
      // Close the window that we opened earlier
      handle.close();
    });
}

const ProfileDropdown: React.FC<{
  instanceName: string | undefined;
  profile: Profile;
  invocationID: string;
}> = ({ instanceName, profile, invocationID }) => {
  const items: MenuProps["items"] = [
    {
      label: "Download Profile",
      key: "download_profile",
      icon: <DownloadOutlined />,
      onClick: () => window.open(getProfileUrl(instanceName, profile), "_self"),
    },
    {
      label: "Open in Perfetto",
      key: "open_in_perfetto",
      icon: <ProjectOutlined rotate={270} />,
      onClick: () => openPerfetto(instanceName, profile, invocationID),
    },
  ];

  return (
    <Dropdown menu={{ items }}>
      <Button>
        <Space>
          Profile
          <DownOutlined />
        </Space>
      </Button>
    </Dropdown>
  );
};

export default ProfileDropdown;
