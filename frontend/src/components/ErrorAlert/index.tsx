import React from 'react';
import { isApolloError } from '@apollo/client';
import PortalAlert from '@/components/PortalAlert';

interface Props {
  error: Error;
}

const ErrorAlert: React.FC<Props> = ({ error }) => (
  <PortalAlert
    type="error"
    showIcon
    message={`Error: ${error.name}`}
    description={isApolloError(error) ? JSON.stringify(error) : error.message}
  />
);

export default ErrorAlert;
