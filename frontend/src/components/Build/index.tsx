import React from 'react';
import linkifyHtml from 'linkify-html';
import { CollapseProps, Descriptions, Space, Typography } from 'antd';
import themeStyles from '@/theme/theme.module.css';
import { FindBuildByUuidQuery } from '@/graphql/__generated__/graphql';
import PortalCard from '@/components/PortalCard';
import PortalAlert from '@/components/PortalAlert';
import BuildStepStatusIcon from '@/components/BuildStepStatusIcon';
import { getFragmentData } from '@/graphql/__generated__';
import {
  BAZEL_INVOCATION_FRAGMENT,
  FULL_BAZEL_INVOCATION_DETAILS, PROBLEM_INFO_FRAGMENT
} from "@/app/bazel-invocations/[invocationID]/index.graphql";
import byResultRank from "@/components/Build/index.helpers";
import { maxBy, xor } from "lodash";
import { BuildStepResultEnum } from "@/components/BuildStepResultTag";
import BazelInvocation from "@/components/BazelInvocation";
import BuildProblems from "@/components/Problems";
import styles from "../AppBar/index.module.css"
import PortalDuration from '../PortalDuration';
import dayjs from 'dayjs';
import Link from 'next/link';

interface Props {
  buildQueryResults: FindBuildByUuidQuery;
  buildStepToDisplayID?: string;
  innerCard?: boolean;
}

const Build: React.FC<Props> = ({ buildQueryResults, buildStepToDisplayID, innerCard }) => {
  const build = buildQueryResults.getBuild
  if (!build) {
    return <></>
  }

  const titleBits: React.ReactNode[] = [
    <span key="build">Build ID: <Typography.Text type="secondary"  className={styles.normalWeight}>{build.buildUUID}</Typography.Text></span>
  ];
  titleBits.push(<span className={styles.copyIcon}>
    <Typography.Text copyable={{text: build.buildUUID ?? "Copy"}}/>
  </span>)
  titleBits.push(<span key="build-url">Build URL: <Link href={build.buildURL} target="_blank" className={styles.normalWeight} ><Typography.Text type="secondary" >{build.buildURL}</Typography.Text> </Link>
  </span>)

  const extraBits: React.ReactNode[] = [
    <Typography.Text code ellipsis className={styles.startedAt}>
      {dayjs(build.timestamp).format('YYYY-MM-DD hh:mm:ss A')}
    </Typography.Text>
  ];

  const invocations = getFragmentData(FULL_BAZEL_INVOCATION_DETAILS, build.invocations);

  const aggregateBuildStepStatus =
    maxBy(
      invocations?.map(invocation => {
        const invocationData = getFragmentData(BAZEL_INVOCATION_FRAGMENT, invocation)
        return invocationData.state.exitCode?.name as BuildStepResultEnum
      }),
      byResultRank,
    ) ?? BuildStepResultEnum.UNKNOWN;

  const envVarItems = build.env.map((env) => {
    return {
      key: env.key,
      label: env.key,
      children: <span dangerouslySetInnerHTML={{ __html: linkifyHtml(env.value) }}></span>
    }
  })

  const items: CollapseProps['items'] = [
    {
      key: '1',
      label: 'This is panel header 1',
      children: <p>1</p>,
    },
    {
      key: '2',
      label: 'This is panel header 2',
      children: <p>2</p>,
    },
    {
      key: '3',
      label: 'This is panel header 3',
      children: <p>3</p>,
    },
  ];

  return (
    <PortalCard
      bordered={false}
      type={innerCard ? 'inner' : undefined}
      icon={<BuildStepStatusIcon status={aggregateBuildStepStatus} />}
      titleBits={titleBits}
      extraBits={extraBits}
    >

      {build.invocations ? (
        <Space direction="vertical" size="middle" className={themeStyles.space}>
          <>
            {
              invocations?.map(invocation => {
                const invocationOverview = getFragmentData(BAZEL_INVOCATION_FRAGMENT, invocation)
                //const problems = invocation.problems.map(p => getFragmentData(PROBLEM_INFO_FRAGMENT, p))
                return (
                  <BazelInvocation
                    key={invocationOverview.invocationID}
                    invocationOverview={invocationOverview}

                    isNestedWithinBuildCard
                  >
                    <BuildProblems
                      problems={[]}
                    />
                  </BazelInvocation>
                );
              })
            }
          </>
        </Space>
      ) : (
        <PortalAlert
          type="success"
          message={<Typography.Title level={5}>No Reported Failures</Typography.Title>}
          description="There is no debug information to display because there are no reported failures with the build"
          showIcon
        />
      )}
    </PortalCard>
  );
};

export default Build;
