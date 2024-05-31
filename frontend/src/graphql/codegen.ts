import { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: [
    '../internal/graphql/schema/scalars.graphql',
    '../internal/graphql/schema/ent.graphql',
    '../internal/graphql/schema/custom.graphql',
  ],
  documents: ['./src/**/*.ts?(x)'],
  generates: {
    './src/graphql/__generated__/': {
      preset: 'client',
      presetConfig: {
        gqlTagName: 'gql',
        fragmentMasking: { unmaskFunctionName: 'getFragmentData' },
        persistedDocuments: true,
      },
    },
  },
};

export default config;
