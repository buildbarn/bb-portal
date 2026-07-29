import {
  DownloadOutlined,
  DownOutlined,
  ProjectOutlined,
} from "@ant-design/icons";
import { Button, Dropdown, type MenuProps, Space } from "antd";
import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { generateFileUrlFromGraphqlFile } from "@/utils/urlGenerator";

const PERFETTO_URL = "https://ui.perfetto.dev";

const fetchProfileFile = async (
  profile: FileDetailsFragment,
): Promise<ArrayBuffer> => {
  const res = await fetch(generateFileUrlFromGraphqlFile(profile));
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

function openPerfetto(profile: FileDetailsFragment, invocationID: string) {
  const handle = window.open(PERFETTO_URL);
  if (!handle) {
    console.error("Failed to open new window for Perfetto UI");
    return;
  }

  Promise.all([fetchProfileFile(profile), waitForPerfettoToLoad(handle)])
    .then((values) => {
      handle.postMessage(
        {
          perfetto: {
            buffer: values[0],
            title: invocationID,
            fileName: profile.filePath,
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
  profile: FileDetailsFragment;
  invocationID: string;
}> = ({ profile, invocationID }) => {
  const items: MenuProps["items"] = [
    {
      label: "Download Profile",
      key: "download_profile",
      icon: <DownloadOutlined />,
      onClick: () =>
        window.open(generateFileUrlFromGraphqlFile(profile), "_self"),
    },
    {
      label: "Open in Perfetto",
      key: "open_in_perfetto",
      icon: <ProjectOutlined rotate={270} />,
      onClick: () => openPerfetto(profile, invocationID),
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
