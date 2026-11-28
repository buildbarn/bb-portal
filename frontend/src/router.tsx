import { createRouter, type RouterHistory } from "@tanstack/react-router";
import { DisabledPage } from "@/components/pages/DisabledPage";
import { ErrorPage } from "@/components/pages/ErrorPage";
import { NotFoundPage } from "@/components/pages/NotFound";
import { NotFoundError } from "@/utils/NotFoundError";
import { routeTree } from "./routeTree.gen";
import {
  parsePortalSearch,
  stringifyPortalSearch,
} from "./schedulers/urlSearch";
import { FeatureDisabledError } from "./utils/featureGuard";

export function createPortalRouter(history?: RouterHistory) {
  return createRouter({
    history,
    routeTree,
    parseSearch: parsePortalSearch,
    stringifySearch: stringifyPortalSearch,
    defaultPreload: "intent",
    defaultNotFoundComponent: () => <NotFoundPage />,
    defaultErrorComponent: ({ error }) => {
      switch (true) {
        case error instanceof FeatureDisabledError:
          return <DisabledPage />;
        case error instanceof NotFoundError:
          return (
            <NotFoundPage
              type={error.type}
              details={error.details}
              showUnauthenticatedMessage={true}
            />
          );
        default:
          return <ErrorPage error={error} />;
      }
    },
  });
}
