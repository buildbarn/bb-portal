import { Button, Space } from "antd";
import { useEffect, useState } from "react";
import { LinkButton } from "@/components/LinkButton";
import { digestFunction_ValueFromJSON } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import {
  type BrowserPageParams,
  BrowserPageSchema,
} from "@/types/BrowserPageParams";
import { BrowserPageType } from "@/types/BrowserPageType";
import { useCheckDataExists } from "@/utils/fetchCasObject";
import { generateBrowserSplat } from "@/utils/urlGenerator";

const COMPARE_KEY = "buildbarn_compare_action";

type Props = {
  params: BrowserPageParams;
  comparing: boolean;
};

const CompareActionButtons: React.FC<Props> = ({ params, comparing }) => {
  const [storedData, setStoredData] = useState<BrowserPageParams | undefined>();

  const { exists } = useCheckDataExists(
    storedData?.instanceName ?? "",
    [
      storedData?.digest ?? {
        hash: "",
        sizeBytes: "",
      },
    ],
    storedData?.digestFunction ?? digestFunction_ValueFromJSON(""),
    storedData !== undefined,
  );

  useEffect(() => {
    const checkSavedData = () => {
      try {
        const rawData = localStorage.getItem(COMPARE_KEY);
        if (!rawData) {
          setStoredData(undefined);
          return;
        }
        const item = BrowserPageSchema.parse({
          ...JSON.parse(rawData),
          browserPageType: BrowserPageType.Action,
        });
        setStoredData(item);
      } catch (error) {
        console.error("Error loading from localStorage", error);
        setStoredData(undefined);
      }
    };

    window.addEventListener("storage", checkSavedData);
    checkSavedData();

    return () => {
      window.removeEventListener("storage", checkSavedData);
    };
  }, []);
  const clearData = () => {
    localStorage.removeItem(COMPARE_KEY);
    setStoredData(undefined);
  };
  const setCompare = () => {
    const comparedAction: Omit<BrowserPageParams, "browserPageType"> & {
      browserPageType?: BrowserPageType;
    } = BrowserPageSchema.parse(params);
    delete comparedAction.browserPageType;
    localStorage.setItem(COMPARE_KEY, JSON.stringify(comparedAction));
    setStoredData(params);
  };

  if (comparing) {
    return (
      <LinkButton to="/browser/$" search={undefined}>
        Stop comparing
      </LinkButton>
    );
  }

  if (!storedData || !exists) {
    return (
      <Button type="default" onClick={setCompare}>
        Compare...
      </Button>
    );
  }
  return (
    <Space.Compact>
      {!(params.digest.hash === storedData.digest.hash) && (
        <LinkButton
          to="/browser/$"
          params={{
            _splat: generateBrowserSplat(
              params.instanceName,
              params.digestFunction,
              params.digest,
              BrowserPageType.Action,
            ),
          }}
          search={{ comparedAction: storedData }}
        >
          Start comparing
        </LinkButton>
      )}
      <Button danger onClick={clearData}>
        Clear compare target
      </Button>
    </Space.Compact>
  );
};

export { CompareActionButtons };
