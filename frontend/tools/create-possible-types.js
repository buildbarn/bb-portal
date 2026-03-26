import fs from "node:fs";
import path from "node:path";
import { buildSchema } from "graphql";

const schemaFiles = [
  "../internal/graphql/schema/scalars.graphql",
  "../internal/graphql/schema/ent.graphql",
  "../internal/graphql/schema/custom.graphql",
];

const schemaString = schemaFiles
  .map((file) => fs.readFileSync(path.resolve(file), "utf8"))
  .join("\n");

const schema = buildSchema(schemaString);
const typeMap = schema.getTypeMap();

const possibleTypes = {};

for (const type of Object.values(typeMap)) {
  possibleTypes[type.name] ??= [];
}

for (const type of Object.values(typeMap)) {
  const interfaces = type.getInterfaces?.() ?? [];
  for (const iface of interfaces) {
    if (possibleTypes[iface.name] !== undefined) {
      possibleTypes[iface.name].push(type.name);
    }
  }
}

fs.writeFileSync(
  "./src/components/ApolloWrapper/possibleTypes.json",
  JSON.stringify(possibleTypes, Object.keys(possibleTypes).sort(), 2),
);
