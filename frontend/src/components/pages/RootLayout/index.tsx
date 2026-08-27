import { TanStackDevtools } from "@tanstack/react-devtools";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtoolsPanel } from "@tanstack/react-query-devtools";
import { Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { useLayoutEffect, useState } from "react";
import { ApolloWrapper } from "@/components/ApolloWrapper";
import { PageWrapper } from "@/components/PageWrapper";
import MessageProvider from "@/context/MessageProvider";
import { ThemeProvider } from "@/context/ThemeProvider";
import { Status } from "@/lib/grpc-client/google/rpc/status";
import { isRetryableGrpcError } from "@/utils/grpcStatus";

export const RootLayout = () => {
  useLayoutEffect(() => {
    document.getElementById("splash")?.remove();
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
      <ThemeProvider>
        <QueryClientProvider client={queryClient}>
          <MessageProvider>
            <PageWrapper>
              <Outlet />
            </PageWrapper>
          </MessageProvider>
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
      </ThemeProvider>
    </ApolloWrapper>
  );
};
