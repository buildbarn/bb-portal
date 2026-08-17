import { ApolloClient, InMemoryCache } from "@apollo/client/core";
import { parse } from "graphql";
import { describe, expect, it } from "vitest";
import { createUnsignedLongHttpLink } from "./unsignedLong";

interface TestQueryData {
  actionCacheStatistics: {
    hits: number;
    saveTimeInMs: bigint;
    sizeInBytes: bigint;
  };
  systemNetworkStats: {
    peakPacketsRecvPerSec: bigint;
  };
}

interface TestQueryVariables {
  value: string;
}

describe("UnsignedLong GraphQL transport", () => {
  it("preserves responses and serializes variables across the uint64 range", async () => {
    let requestVariables: Record<string, unknown> | undefined;
    const fetcher: typeof fetch = async (_input, init) => {
      if (typeof init?.body !== "string") {
        throw new Error("Expected a JSON GraphQL request body");
      }
      const request = JSON.parse(init.body) as {
        variables?: Record<string, unknown>;
      };
      requestVariables = request.variables;

      return new Response(
        `{
          "data": {
            "actionCacheStatistics": {
              "__typename": "ActionCacheStatistics",
              "hits": 7,
              "saveTimeInMs": 42,
              "sizeInBytes": 9007199254740993
            },
            "systemNetworkStats": {
              "__typename": "SystemNetworkStats",
              "peakPacketsRecvPerSec": 18446744073709551615
            }
          }
        }`,
        {
          headers: { "content-type": "application/graphql-response+json" },
          status: 200,
        },
      );
    };
    const client = new ApolloClient({
      cache: new InMemoryCache(),
      link: createUnsignedLongHttpLink({ fetch: fetcher }),
    });

    const result = await client.query<TestQueryData, TestQueryVariables>({
      query: parse(`
        query UnsignedLongTransport($value: UnsignedLong!) {
          actionCacheStatistics(value: $value) {
            hits
            saveTimeInMs
            sizeInBytes
          }
          systemNetworkStats {
            peakPacketsRecvPerSec
          }
        }
      `),
      variables: { value: "18446744073709551615" },
    });

    expect(requestVariables).toEqual({ value: "18446744073709551615" });
    expect(result.data).toBeDefined();
    if (result.data !== undefined) {
      expect(result.data.actionCacheStatistics.hits).toBe(7);
      expect(result.data.actionCacheStatistics.saveTimeInMs).toBe(42n);
      expect(result.data.actionCacheStatistics.sizeInBytes).toBe(
        9007199254740993n,
      );
      expect(result.data.systemNetworkStats.peakPacketsRecvPerSec).toBe(
        18446744073709551615n,
      );
    }
  });
});
