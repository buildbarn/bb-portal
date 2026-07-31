import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: [
    "../internal/graphql/schema/scalars.graphql",
    "../internal/graphql/schema/ent.graphql",
    "../internal/graphql/schema/custom.graphql",
  ],
  documents: ["./src/**/*.ts?(x)"],
  generates: {
    "./src/graphql/__generated__/": {
      preset: "client",
      presetConfig: {
        gqlTagName: "gql",
        fragmentMasking: { unmaskFunctionName: "getFragmentData" },
        persistedDocuments: true,
      },
      config: {
        scalars: {
          UnsignedLong: {
            input: "string",
            output: "bigint",
          },
        },
        useTypeImports: true,
      },
    },
    "./src/graphql/__generated__/zod.ts": {
      plugins: ["typescript-validation-schema"],
      config: {
        schema: "zodv4",
        scalarSchemas: {
          Time: "z.iso.datetime({ offset: true })",
          UnsignedLong:
            'z.string().regex(/^\\d+$/, "Expected an unsigned integer")',
          UUID: "z.uuid()",
        },
        importFrom: "./graphql",
        notAllowEmptyString: true,
        useEnumTypeAsDefaultValue: true,
      },
    },
  },
};

export default config;
