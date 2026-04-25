import { TanStackDevtools } from "@tanstack/react-devtools";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtoolsPanel } from "@tanstack/react-query-devtools";
import { Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { ConfigProvider, Layout } from "antd";
import { useCallback, useLayoutEffect, useState } from "react";
import { ApolloWrapper } from "@/components/ApolloWrapper";
import { PageWrapper } from "@/components/PageWrapper";
import { Status } from "@/lib/grpc-client/google/rpc/status";
import dark from "@/theme/dark";
import light from "@/theme/light";
import { isRetryableGrpcError } from "@/utils/grpcStatus";
import styles from "./index.module.css";

const PREFERS_DARK_KEY = "prefers-dark";

function getTheme(): "dark" | "light" {
  return document.documentElement.getAttribute("data-theme") === "dark"
    ? "dark"
    : "light";
}

function setTheme(theme: "dark" | "light") {
  window.localStorage.setItem(
    PREFERS_DARK_KEY,
    theme === "dark" ? "true" : "false",
  );
  document.documentElement.setAttribute("data-theme", theme);
}

export const RootLayout = () => {
  // Inner theme is a meaningless state used to force rerender when we modify the html tag.
  const [innerTheme, setInnerTheme] = useState(() => getTheme());

  useLayoutEffect(() => {
    document.getElementById("splash")?.remove();
  }, []);

  const toggleTheme = useCallback(() => {
    const opposite = getTheme() === "dark" ? "light" : "dark";
    setTheme(opposite);
    setInnerTheme(opposite);
  }, []);

  const [queryClient] = useState(
    () =>
      new QueryClient({
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
      }),
  );

  return (
    <ApolloWrapper>
      <ConfigProvider theme={innerTheme === "dark" ? dark : light}>
        <QueryClientProvider client={queryClient}>
          <Layout className={styles.layout}>
            <PageWrapper
              toggleTheme={toggleTheme}
              prefersDark={innerTheme === "dark"}
            >
              <Outlet />
            </PageWrapper>
          </Layout>
          {/* Devtools for Tanstack components. Automatically removed for prod builds */}
          <TanStackDevtools
            plugins={[
              {
                name: "TanStack Query",
                render: <ReactQueryDevtoolsPanel />,
              },
              {
                name: "TanStack Router",
                render: <TanStackRouterDevtoolsPanel />,
              },
            ]}
          />
        </QueryClientProvider>
      </ConfigProvider>
    </ApolloWrapper>
  );
};
