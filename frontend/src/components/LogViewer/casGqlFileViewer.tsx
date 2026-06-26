import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { CasViewer } from "./casViewer";

interface Props {
  file: FileDetailsFragment;
  title: string;
  fileName: string;
}

export const CasGqlFileViewer: React.FC<Props> = ({
  file,
  title,
  fileName,
}) => {
  return (
    <CasViewer
      instanceName={file.digest.rev2InstanceName}
      digestFunction={file.digest.digestFunction}
      hash={file.digest.hash}
      sizeBytes={file.digest.sizeBytes}
      title={title}
      fileName={fileName}
    />
  );
};
