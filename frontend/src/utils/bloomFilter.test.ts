import { Digest } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import { FileSystemAccessProfile } from "@/lib/grpc-client/buildbarn/fsac/fsac";
import { expect, test } from "vitest";
import {
  PATH_HASH_BASE_HASH,
  PathHashes,
  containsPathHashes,
  generateFileSystemReferenceQueryParams,
  readBloomFilter,
} from "./bloomFilter";

test("bloomFilterReader", () => {
  expect(() => readBloomFilter(FileSystemAccessProfile.create())).toThrowError(
    "Bloom filter is empty",
  );

  expect(() =>
    readBloomFilter({
      bloomFilter: Uint8Array.from([0x01]),
      bloomFilterHashFunctions: 123,
    }),
  ).toThrowError("Bloom filter has zero bits");

  expect(() =>
    readBloomFilter({
      bloomFilter: Uint8Array.from([0x12, 0x00]),
      bloomFilterHashFunctions: 123,
    }),
  ).toThrowError("Bloom filter's trailing byte is not properly padded");
});

test("containsPathHashes", () => {
  const bloomFilterReader = readBloomFilter({
    bloomFilter: Uint8Array.from([
      0x1d, 0xb2, 0x43, 0xf1, 0x61, 0xfa, 0x18, 0x3f,
    ]),
    bloomFilterHashFunctions: 11,
  });

  expect(
    containsPathHashes(
      bloomFilterReader,
      new PathHashes().appendComponent("dir"),
    ),
  ).toBe(true);

  expect(
    containsPathHashes(
      bloomFilterReader,
      new PathHashes().appendComponent("file"),
    ),
  ).toBe(true);

  expect(
    containsPathHashes(
      bloomFilterReader,
      new PathHashes().appendComponent("dir").appendComponent("file"),
    ),
  ).toBe(true);

  expect(
    containsPathHashes(
      bloomFilterReader,
      new PathHashes().appendComponent("nonexistent"),
    ),
  ).toBe(false);
});

test("generateFileSystemReferenceQueryParams", () => {
  expect(generateFileSystemReferenceQueryParams(undefined)).toBeUndefined();

  expect(
    generateFileSystemReferenceQueryParams({
      digest: Digest.create({
        hash: "01234",
        sizeBytes: "999",
      }),
      pathHashesBaseHash: "56789",
    }),
  ).toEqual({
    fileSystemAccessProfile:
      "%7B%22digest%22%3A%7B%22hash%22%3A%2201234%22%2C%22sizeBytes%22%3A%22999%22%7D%2C%22pathHashesBaseHash%22%3A%2256789%22%7D",
  });
});

test("PathHashes", () => {
  expect(new PathHashes().baseHash).toEqual(BigInt(PATH_HASH_BASE_HASH));

  expect(new PathHashes(BigInt("123456789")).baseHash).toEqual(
    BigInt("123456789"),
  );
});
