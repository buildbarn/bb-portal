import {CheckCircleOutlined, CloseCircleFilled, QuestionCircleOutlined, StopOutlined} from '@ant-design/icons';
import React from 'react';
import {BuildStepStatus} from '@/graphql/__generated__/graphql';
import themeStyles from '@/theme/theme.module.css';
import {BuildStepResultEnum} from "@/components/BuildStepResultTag";

interface Props {
  status: BuildStepResultEnum;
}

const BuildStepStatusIcon: React.FC<Props> = ({ status }) => {
  switch (status) {
    case BuildStepResultEnum.TESTS_FAILED  || BuildStepResultEnum.BUILD_FAILURE || BuildStepResultEnum.PARSING_FAILURE : {
      return <CloseCircleFilled className={themeStyles.colorFailure} />;
    }
    case BuildStepResultEnum.SUCCESS: {
      return <CheckCircleOutlined className={themeStyles.colorSuccess} />;
    }
    case BuildStepResultEnum.ABORTED || BuildStepResultEnum.INTERRUPTED : {
      return <StopOutlined />;
    }
    default: {
      return <QuestionCircleOutlined />;
    }
  }
};

export default BuildStepStatusIcon;
