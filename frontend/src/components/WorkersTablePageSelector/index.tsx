import type { PaginationInfo } from "@/lib/grpc-client/buildbarn/buildqueuestate/buildqueuestate";
import { BackwardFilled, CaretRightFilled } from "@ant-design/icons";
import { Button, Space, Typography } from "antd";
import type React from "react";

interface Props {
  paginationInfo: PaginationInfo;
  goToFirstPage: () => void;
  goToNextPage: () => void;
  pageSize: number;
}

const WorkersTablePageSelector: React.FC<Props> = ({
  paginationInfo,
  goToFirstPage,
  goToNextPage,
  pageSize,
}) => {
  return (
    <Space>
      <Button
        type="primary"
        icon={<BackwardFilled />}
        onClick={goToFirstPage}
      />
      <Typography.Text>
        {`Showing workers [${paginationInfo.startIndex}, ${Math.min(paginationInfo.startIndex + pageSize, paginationInfo.totalEntries)}) of ${paginationInfo.totalEntries} in total.`}
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
