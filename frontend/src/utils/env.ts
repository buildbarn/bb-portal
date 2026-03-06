import { PortalFrontendConfiguration } from "@/lib/grpc-client/portal/frontend/frontend";

export const env: PortalFrontendConfiguration = window["__env__"]

if (!env) {
  throw new Error("Missing frontend configuration: The frontend needs to be served through the backend to receive the configuration. Without it the frontend will not work.")
}
