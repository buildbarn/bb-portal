import { env } from "@/utils/env";

export const generatePageTitle = (elements: (string | undefined)[]): string => {
  return [
    ...elements.filter((elem) => !!elem),
    `${env.companyName ? `${env.companyName} ` : ""}Buildbarn Portal`,
  ].join(" - ");
};
