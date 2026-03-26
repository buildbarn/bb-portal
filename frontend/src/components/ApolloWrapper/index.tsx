import type React from 'react';
import { 
  ApolloClient,
  InMemoryCache,
  HttpLink 
} from '@apollo/client';
import { ApolloProvider } from '@apollo/client/react';
import possibleTypes from './possibleTypes.json';

export const apolloClient = new ApolloClient({
  link: new HttpLink({
    uri: "/graphql",
  }),
  cache: new InMemoryCache({
    possibleTypes
  }),
});


export const ApolloWrapper = ({ children }: React.PropsWithChildren) => {
  return <ApolloProvider client={apolloClient}>{children}</ApolloProvider>;
};
