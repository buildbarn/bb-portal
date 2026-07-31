import { ApolloClient, InMemoryCache } from "@apollo/client/core";
import { ApolloProvider } from "@apollo/client/react";
import type React from "react";
import { createUnsignedLongHttpLink } from "@/graphql/unsignedLong";
import possibleTypes from "./possibleTypes.json";

export const apolloClient = new ApolloClient({
  link: createUnsignedLongHttpLink(),
  cache: new InMemoryCache({
    possibleTypes,
  }),
});

export const ApolloWrapper = ({ children }: React.PropsWithChildren) => {
  return <ApolloProvider client={apolloClient}>{children}</ApolloProvider>;
};
