import type { PortalFrontendConfiguration } from "@/lib/grpc-client/portal/frontend/frontend";

// biome-ignore lint/suspicious/noExplicitAny: We have no type information, and it is checked by the backend anyways.
export const env: PortalFrontendConfiguration = (window as any).__env__;

if (!env) {
  throw new Error(
    "Missing frontend configuration: The frontend needs to be served through the backend to receive the configuration. Without it the frontend will not work.",
  );
}
