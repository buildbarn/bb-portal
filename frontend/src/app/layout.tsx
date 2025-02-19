'use client';

import '@/app/globals.css';
import { ConfigProvider, Layout } from 'antd';
import React, { useCallback, useEffect, useState } from 'react';
import { ItemType } from 'antd/lib/menu/hooks/useItems';
import styles from '@/app/layout.module.css';
import RootStyleRegistry from '@/components/RootStyleRegistry';
import AppBar, { SetExtraAppBarMenuItemsContext } from '@/components/AppBar';
import dark from '@/theme/dark';
import light from '@/theme/light';
import Dynamic from '@/components/Dynamic';
import { ApolloWrapper } from '@/components/ApolloWrapper';
import parseStringBoolean from '@/utils/storage';
import { PublicEnvScript } from 'next-runtime-env';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import GrpcClientsProvider from '@/context/GrpcClientsProvider';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

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

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const [prefersDark, setPrefersDark] = useState<boolean>(calculatePrefersDark());
  const theme = prefersDark ? dark : light;

  useEffect(() => {
    const val = calculatePrefersDark();
    setPrefersDark(val);
  }, [setPrefersDark]);

  const toggleTheme = useCallback(() => {
    const opposite = !prefersDark;
    savePrefersDark(opposite);
    setPrefersDark(opposite);
  }, [prefersDark, setPrefersDark]);

  const [extraAppBarMenuItems, setExtraAppBarMenuItems] = useState<ItemType[]>([]);

  const queryClient = new QueryClient()

  return (
    <>
      <title>Buildbarn Portal</title>
      <html lang="en" className={styles.html}>
        <head>
          <PublicEnvScript />
        </head>
        <body className={styles.body}>
          <div style={{ display: "flex" }}>
            <ApolloWrapper>
              <RootStyleRegistry>
                <ConfigProvider theme={theme}>
                  <QueryClientProvider client={queryClient}>
                    <GrpcClientsProvider>
                      <Dynamic>
                        <Layout className={styles.layout}>
                          <AppBar
                            toggleTheme={toggleTheme}
                            prefersDark={prefersDark}
                            extraMenuItems={extraAppBarMenuItems}
                          />
                          <SetExtraAppBarMenuItemsContext.Provider value={setExtraAppBarMenuItems}>
                            {children}
                          </SetExtraAppBarMenuItemsContext.Provider>
                        </Layout>
                      </Dynamic>
                    </GrpcClientsProvider>
                    {/* Adds devtools. Is automatically removed for production builds. */}
                    <ReactQueryDevtools initialIsOpen={false} />
                  </QueryClientProvider>
                </ConfigProvider>
              </RootStyleRegistry>
            </ApolloWrapper>
          </div>
        </body>
      </html>
    </>
  );
}
