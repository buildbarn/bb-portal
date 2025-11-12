import { Collapse, Descriptions, Space, Table, Typography } from "antd";
import type {
  AuthenticatedUserNodeFragmentFragment,
  BazelInvocationNodeFragment,
} from "@/graphql/__generated__/graphql";
import themeStyles from "@/theme/theme.module.css";
import getColumns from "./columns";

interface Props {
  authenticatedUser: AuthenticatedUserNodeFragmentFragment | null | undefined;
}

const UserView: React.FC<Props> = ({ authenticatedUser }) => {
  const invocations =
    authenticatedUser?.bazelInvocations.edges?.map((node) => node?.node) || [];
  const userInfo = authenticatedUser?.userInfo || {};

  return (
    <Space direction="vertical" className={themeStyles.space}>
      {Object.keys(userInfo).length > 0 && (
        <Collapse
          bordered={false}
          items={[
            {
              key: 0,
              label: <Typography.Text strong>User information</Typography.Text>,
              children: (
                <Descriptions column={1} bordered size="small">
                  {Object.keys(userInfo).map((value) => {
                    return (
                      <Descriptions.Item key={value} label={value}>
                        {userInfo[value]}
                      </Descriptions.Item>
                    );
                  })}
                </Descriptions>
              ),
            },
          ]}
        />
      )}
      <Table
        dataSource={invocations as BazelInvocationNodeFragment[]}
        columns={getColumns()}
        size="small"
        pagination={false}
      />
    </Space>
  );
};

export default UserView;
