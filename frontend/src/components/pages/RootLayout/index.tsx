import { ApolloWrapper } from '@/components/ApolloWrapper';
import { ConfigProvider, Layout } from 'antd';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import GrpcClientsProvider from '@/context/GrpcClientsProvider';
import AppBar from '@/components/AppBar';
import { TanStackDevtools } from '@tanstack/react-devtools';
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools';
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';
import { useCallback, useEffect, useLayoutEffect, useState } from 'react';
import parseStringBoolean from '@/utils/storage';
import dark from '@/theme/dark';
import light from '@/theme/light';
import { isRetryableGrpcError } from '@/utils/grpcStatus';
import { Status } from '@/lib/grpc-client/google/rpc/status';
import { Outlet } from '@tanstack/react-router';
import styles from './index.module.css';

const PREFERS_DARK_KEY = 'prefers-dark';

function windowDefined(): boolean {
  return typeof window !== 'undefined';
}

function calculatePrefersDark(): boolean {
  if (windowDefined()) {
    const prefersDark = window.localStorage.getItem(PREFERS_DARK_KEY);
    if (prefersDark) {
      return parseStringBoolean(prefersDark);
    }
  }
  return browserPrefersDark();
}

function browserPrefersDark(): boolean {
  if (windowDefined()) {
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  }
  return false;
}

function savePrefersDark(prefersDark: boolean): void {
  if (windowDefined()) {
    window.localStorage.setItem(PREFERS_DARK_KEY, prefersDark ? 'true' : 'false');
    window.localStorage.setItem('graphiql:theme', prefersDark ? 'dark' : 'light');
  }
}

export const RootLayout = () => {
  const [prefersDark, setPrefersDark] = useState<boolean>(calculatePrefersDark());
  const theme = prefersDark ? dark : light;

  useLayoutEffect(() => {
    document.getElementById('splash')?.remove();
  }, [])

  useEffect(() => {
    const val = calculatePrefersDark();
    setPrefersDark(val);
  }, [setPrefersDark]);

  const toggleTheme = useCallback(() => {
    const opposite = !prefersDark;
    savePrefersDark(opposite);
    setPrefersDark(opposite);
  }, [prefersDark, setPrefersDark]);

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount: number, error: Error) => {
          if (failureCount >= 3) {
            return false;
          }
          return isRetryableGrpcError(Status.fromJSON(error));
        },
      },
    },
  });
  return <ApolloWrapper>
    <ConfigProvider theme={theme}>
      <QueryClientProvider client={queryClient}>
        <GrpcClientsProvider>
          <Layout className={styles.layout}>
            <AppBar
              toggleTheme={toggleTheme}
              prefersDark={prefersDark}
            />
            <Outlet />
          </Layout>
        </GrpcClientsProvider>
        {/* Devtools for Tanstack components. Automatically removed for prod builds */}
        <TanStackDevtools
          plugins={[
            {
              name: 'TanStack Query',
              render: <ReactQueryDevtoolsPanel />,
            },
            {
              name: 'TanStack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
          ]}
        />
      </QueryClientProvider>
    </ConfigProvider>
  </ApolloWrapper>
}
