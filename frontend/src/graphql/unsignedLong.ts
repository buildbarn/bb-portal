import { ApolloLink, HttpLink, Observable } from "@apollo/client/core";
import { type DocumentNode, visit } from "graphql";
import { isLosslessNumber, parse as parseLosslessJSON } from "lossless-json";

const unsignedLongFieldNames = new Set([
  "sizeInBytes",
  "saveTimeInMs",
  "loadTimeInMs",
  "cacheCheckSemaphoreWaitTimeInMs",
  "bytesSent",
  "bytesRecv",
  "packetsSent",
  "packetsRecv",
  "peakBytesSentPerSec",
  "peakBytesRecvPerSec",
  "peakPacketsSentPerSec",
  "peakPacketsRecvPerSec",
]);

const losslessNumberKey = "__bbPortalLosslessNumber";

interface LosslessNumberEnvelope {
  [losslessNumberKey]: string;
}

interface UnsignedLongHttpLinkOptions {
  uri?: string;
  fetch?: typeof globalThis.fetch;
}

const isObject = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null;

const wrapLosslessNumbers = (value: unknown): unknown => {
  if (isLosslessNumber(value)) {
    return { [losslessNumberKey]: value.toString() };
  }
  if (Array.isArray(value)) {
    return value.map(wrapLosslessNumbers);
  }
  if (isObject(value)) {
    return Object.fromEntries(
      Object.entries(value).map(([key, child]) => [
        key,
        wrapLosslessNumbers(child),
      ]),
    );
  }
  return value;
};

const isLosslessNumberEnvelope = (
  value: unknown,
): value is LosslessNumberEnvelope =>
  isObject(value) &&
  Object.keys(value).length === 1 &&
  typeof value[losslessNumberKey] === "string";

const unsignedLongResponseKeys = (document: DocumentNode): Set<string> => {
  const responseKeys = new Set<string>();
  visit(document, {
    Field(node) {
      if (unsignedLongFieldNames.has(node.name.value)) {
        responseKeys.add(node.alias?.value ?? node.name.value);
      }
    },
  });
  return responseKeys;
};

const reviveResponseNumbers = (
  value: unknown,
  unsignedLongKeys: ReadonlySet<string>,
  fieldName?: string,
): unknown => {
  if (isLosslessNumberEnvelope(value)) {
    return fieldName !== undefined && unsignedLongKeys.has(fieldName)
      ? BigInt(value[losslessNumberKey])
      : Number(value[losslessNumberKey]);
  }
  if (Array.isArray(value)) {
    return value.map((child) =>
      reviveResponseNumbers(child, unsignedLongKeys, fieldName),
    );
  }
  if (isObject(value)) {
    return Object.fromEntries(
      Object.entries(value).map(([key, child]) => [
        key,
        reviveResponseNumbers(child, unsignedLongKeys, key),
      ]),
    );
  }
  return value;
};

const responseWithBody = (response: Response, body: BodyInit | null) => {
  const headers = new Headers(response.headers);
  headers.delete("content-encoding");
  headers.delete("content-length");
  headers.delete("transfer-encoding");
  return new Response(body, {
    headers,
    status: response.status,
    statusText: response.statusText,
  });
};

const isJSONResponse = (response: Response) => {
  const mediaType = response.headers
    .get("content-type")
    ?.split(";", 1)[0]
    .trim()
    .toLowerCase();
  return mediaType === "application/json" || mediaType?.endsWith("+json");
};

const losslessJSONFetch = (fetcher: typeof globalThis.fetch) => {
  const wrappedFetch: typeof globalThis.fetch = async (input, init) => {
    const response = await fetcher(input, init);
    if (!isJSONResponse(response)) {
      return response;
    }

    const body = await response.text();
    if (body.length === 0) {
      return responseWithBody(response, null);
    }

    const parsedBody = parseLosslessJSON(body);
    return responseWithBody(
      response,
      JSON.stringify(wrapLosslessNumbers(parsedBody)),
    );
  };
  return wrappedFetch;
};

const unsignedLongScalarLink = new ApolloLink((operation, forward) => {
  const responseKeys = unsignedLongResponseKeys(operation.query);

  return new Observable((observer) => {
    const subscription = forward(operation).subscribe({
      next: (result) => {
        observer.next(
          reviveResponseNumbers(result, responseKeys) as typeof result,
        );
      },
      error: (error: unknown) => observer.error(error),
      complete: () => observer.complete(),
    });
    return () => subscription.unsubscribe();
  });
});

export const createUnsignedLongHttpLink = ({
  uri = "/graphql",
  fetch: fetcher = globalThis.fetch,
}: UnsignedLongHttpLinkOptions = {}) =>
  ApolloLink.from([
    unsignedLongScalarLink,
    new HttpLink({
      uri,
      fetch: losslessJSONFetch(fetcher),
    }),
  ]);
