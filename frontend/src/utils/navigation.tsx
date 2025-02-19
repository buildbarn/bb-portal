import type { AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime";

const updateQueryParam = (
  router: AppRouterInstance,
  currentSearchParams: URLSearchParams,
  pathname: string,
  param: string,
  value?: string,
) => {
  if (!value) {
    currentSearchParams.delete(param);
  } else {
    currentSearchParams.set(param, value);
  }

  const search = currentSearchParams.toString();
  const query = search ? `?${search}` : "";
  router.push(`${pathname}${query}`);
};

export default updateQueryParam;
