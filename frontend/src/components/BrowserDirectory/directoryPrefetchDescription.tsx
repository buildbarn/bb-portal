import { Typography } from "antd";
import themeStyles from "@/theme/theme.module.css";

interface Params {
  prefetchDataExists: boolean;
}

export const DirectoryPrefetchDescription: React.FC<Params> = ({
  prefetchDataExists,
}) => {
  return (
    <Typography.Text>
      {prefetchDataExists && (
        <>
          <strong>Note:</strong>{" "}
          <span className={themeStyles.colorSuccess}>Green</span> and{" "}
          <span className={themeStyles.colorFailure}>
            <s>red</s>
          </span>{" "}
          filenames indicate which files and directories will be prefetched the
          next time a similar action executes. Though it is representative of
          what is actually accessed by the action, it may contain false
          positives and negatives.
        </>
      )}
      <p>
        To see more information about a directory, hold shift, ctrl, or meta key
        while clicking on the directory.
      </p>
    </Typography.Text>
  );
};
