import {ActionOutputStatus, ActionProblem, BlobReference, GetActionProblemQuery} from "@/graphql/__generated__/graphql";
import {Empty, Space, Spin, Tabs} from "antd";
import React, {useEffect} from "react";
import {NetworkStatus, useQuery} from "@apollo/client";
import ErrorAlert from "@/components/ErrorAlert";
import themeStyles from '@/theme/theme.module.css';
import styles from './index.module.css';
import {GET_ACTION_PROBLEM} from "@/app/bazel-invocations/[invocationID]/index.graphql";
import LogOutput from "@/components/Problems/LogOutput";

const { TabPane } = Tabs;

function isActionProblem(node: any): node is ActionProblem {
  // eslint-disable-next-line no-underscore-dangle
  return node?.__typename === 'ActionProblem';
}

function getActionProblem(node: any): ActionProblem | undefined {
  if (isActionProblem(node)) {
    return node;
  }
  return undefined;
}

function isStillProcessing(actionProblem: ActionProblem): boolean {
  const { stderr, stdout } = actionProblem;
  const stderrProcessing = stderr?.availabilityStatus === ActionOutputStatus.Processing;
  const stdoutProcessing = stdout?.availabilityStatus === ActionOutputStatus.Processing;
  return stderrProcessing || stdoutProcessing;
}

const ActionOutputPanel: React.FC<{ blobReference: BlobReference; instanceName: string | undefined; spin: boolean }> = ({ blobReference, instanceName, spin }) => {
  return (
    <Space direction="vertical" size="middle" className={themeStyles.space}>
      {/* Display spin behind the actions, making UI stable when query is being executed. */}
      <div className={themeStyles.flex}>
        {spin && <Spin />}
      </div>
      <LogOutput blobReference={blobReference} instanceName={instanceName} />
    </Space>
  );
};

const ActionProblemPanel: React.FC<{ actionProblem: ActionProblem; instanceName: string | undefined; spin: boolean }> = ({ actionProblem, instanceName, spin }) => {
  const empty = <Empty description="No action output." />;
  const { stderr, stdout } = actionProblem;

  if (!stderr && !stdout) {
    return empty;
  }
  const stderrPanel = stderr && <ActionOutputPanel blobReference={stderr} instanceName={instanceName} spin={spin} />;
  const stdoutPanel = stdout && <ActionOutputPanel blobReference={stdout} instanceName={instanceName} spin={spin} />;

  if (stderr && !stdout) {
    return stderrPanel ?? null;
  }

  if (!stderr && stdout) {
    return stdoutPanel ?? null;
  }

  // stderr && stdout
  return (
    <Tabs defaultActiveKey="stderr">
      <TabPane tab="stderr" key="stderr">
        {stderrPanel}
      </TabPane>
      <TabPane tab="stdout" key="stdout" className={styles.TabAnimationSafari13Fix}>
        {stdoutPanel}
      </TabPane>
    </Tabs>
  );
};

interface Props {
  id: string;
  instanceName: string | undefined;
}


const ActionProblemContainer: React.FC<Props> = ({ id, instanceName }) => {
  const { data, error, loading, stopPolling, networkStatus } = useQuery<GetActionProblemQuery>(GET_ACTION_PROBLEM, {
    variables: { id },
    fetchPolicy: 'cache-and-network',
    pollInterval: 5000,
    notifyOnNetworkStatusChange: true,
  });

  const node = data?.node;
  const actionProblem = getActionProblem(node);
  const stop = actionProblem !== undefined && !isStillProcessing(actionProblem);

  useEffect(() => {
    if (stop) {
      stopPolling();
    }
  }, [stop, stopPolling]);

  if (loading && networkStatus !== NetworkStatus.poll) {
    return <Spin />;
  }
  if (error) {
    return <ErrorAlert error={new Error('Failed to get action problem')} />;
  }
  if (actionProblem === undefined) {
    return <ErrorAlert error={new Error('Expected action problem but server returned something else')} />;
  }

  return <ActionProblemPanel actionProblem={actionProblem} instanceName={instanceName} spin={networkStatus === NetworkStatus.poll} />;
};

export default ActionProblemContainer;
