'use client';
import React from 'react';
import { ApolloLink, HttpLink } from '@apollo/client';
import {
  NextSSRApolloClient,
  ApolloNextAppProvider,
  NextSSRInMemoryCache,
  SSRMultipartLink,
} from '@apollo/experimental-nextjs-app-support/ssr';
import possibleTypes from './possibleTypes.json';
import { env } from 'next-runtime-env';


export const makeClient = () => {
  const httpLink = new HttpLink({
    uri: `${env('NEXT_PUBLIC_BES_BACKEND_URL')}/graphql`,
    fetchOptions: { cache: "no-store" },
  });

  return new NextSSRApolloClient({
    cache: new NextSSRInMemoryCache({
      possibleTypes
    }),
    connectToDevTools: true,
    link:
      typeof window === 'undefined'
        ? ApolloLink.from([
          new SSRMultipartLink({
            stripDefer: true,
          }),
          httpLink,
        ])
        : httpLink,
  });
}

export const ApolloWrapper = ({ children }: React.PropsWithChildren) => {
  return <ApolloNextAppProvider makeClient={makeClient}>{children}</ApolloNextAppProvider>;
}
