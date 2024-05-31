'use client';

import React from 'react';
import { Skeleton, Typography } from 'antd';
import { useQuery } from '@apollo/client';
import { ExclamationCircleOutlined, LoadingOutlined, QuestionCircleFilled, RocketOutlined } from '@ant-design/icons';
import {FindBuildByUuidQuery } from '@/graphql/__generated__/graphql';
import PortalCard from '@/components/PortalCard';
import PortalAlert from '@/components/PortalAlert';
import {FIND_BUILD_BY_UUID_QUERY} from "@/app/builds/[buildUUID]/[[...slugs]]/index.graphql";
import Build from "@/components/Build";

interface StatusProps {
  buildUUID: string;
}

const Loading: React.FC<StatusProps> = ({ buildUUID }) => {
  return (
    <PortalCard
      bordered={false}
      icon={<LoadingOutlined />}
      titleBits={[<span key="title">Loading Build <Typography.Text code>{buildUUID}</Typography.Text>...</span>]}
    >
      <Skeleton active paragraph={{ rows: 6 }} />
    </PortalCard>
  );
};

const Error: React.FC<StatusProps> = ({ buildUUID }) => {
  return (
    <PortalCard
      bordered={false}
      icon={<ExclamationCircleOutlined />}
      titleBits={[<span key="title">Build <Typography.Text code>{buildUUID}</Typography.Text></span>]}
    >
      <PortalAlert
        icon={<ExclamationCircleOutlined />}
        message={<Typography.Title level={5}>Error Retrieving Build Information</Typography.Title>}
        description="Buildbarn Portal experienced an error in fetching data for the build"
        type="warning"
        showIcon
      />
    </PortalCard>
  );
};

const Waiting: React.FC<StatusProps> = ({ buildUUID }) => {
  return (
    <PortalCard
      bordered={false}
      icon={<RocketOutlined rotate={20} />}
      titleBits={[<span key="title">Build <Typography.Text code>{buildUUID}</Typography.Text></span>]}
    >
      <PortalAlert
        icon={<QuestionCircleFilled />}
        message={<Typography.Title level={5}>Awaiting Build Information</Typography.Title>}
        description="Buildbarn Portal is awaiting data for the build"
        type="info"
        showIcon
      />
    </PortalCard>
  );
};

export default function Page({ params }: { params: { buildUUID: string } }) {
  const { loading, error, data } = useQuery<FindBuildByUuidQuery>(FIND_BUILD_BY_UUID_QUERY, {
    variables: { uuid: params.buildUUID },
    fetchPolicy: 'no-cache',
    pollInterval: 60000,
  });

  if (loading) {
    return <Loading buildUUID={params.buildUUID} />;
  } else if (error) {
    console.error(error);
    return <Error buildUUID={params.buildUUID} />;
  } else if (!data?.getBuild) {
    return <Waiting buildUUID={params.buildUUID} />;
  }

  return (
    <Build
      buildQueryResults={data}
    />
  );
}
