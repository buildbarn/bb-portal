import { Typography } from "antd";
import { env } from "@/utils/env";
import { Link } from '@tanstack/react-router';

const AppBarTitle = () => {
  return (
    <Link to="/">
      <Typography.Title level={3}>
        {`${env.companyName ? `${env.companyName} ` : ''}Buildbarn Portal`}
      </Typography.Title>
    </Link>
  );
};

export default AppBarTitle;
