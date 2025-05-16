"use client";
import { Typography } from "antd";
import { env } from "next-runtime-env";
import Link from "next/link";

const AppBarTitle = () => {
  return (
    <Link href="/">
      <Typography.Title level={3}>
        {env("NEXT_PUBLIC_COMPANY_NAME")} Buildbarn Portal
      </Typography.Title>
    </Link>
  );
};

export default AppBarTitle;
