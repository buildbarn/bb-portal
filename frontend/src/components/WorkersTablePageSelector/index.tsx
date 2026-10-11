import { BackwardFilled, CaretRightFilled } from "@ant-design/icons";
import { Button, Space, Typography } from "antd";
import type React from "react";
import type { PaginationInfo } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";

interface Props {
  paginationInfo: PaginationInfo;
  goToFirstPage: () => void;
  goToNextPage: () => void;
  pageSize: number;
  noun?: string;
}

const WorkersTablePageSelector: React.FC<Props> = ({
  paginationInfo,
  goToFirstPage,
  goToNextPage,
  pageSize,
  noun = "workers",
}) => {
  return (
    <Space>
      <Button
        type="primary"
        icon={<BackwardFilled />}
        onClick={goToFirstPage}
      />
      <Typography.Text>
        {`Showing ${noun} [${paginationInfo.startIndex}, ${Math.min(paginationInfo.startIndex + pageSize, paginationInfo.totalEntries)}) of ${paginationInfo.totalEntries} in total.`}
      </Typography.Text>
      <Button
        type="primary"
        icon={<CaretRightFilled />}
        onClick={goToNextPage}
      />
    </Space>
  );
};

export default WorkersTablePageSelector;
