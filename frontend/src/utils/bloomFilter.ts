import type { ParsedUrlQueryInput } from "node:querystring";
import type { FileSystemAccessProfile } from "@/lib/grpc-client/buildbarn/fsac/fsac";
import { FileSystemAccessProfileReference } from "@/lib/grpc-client/buildbarn/query/query";

export interface BloomFilterReader {
  bloomFilter: Uint8Array;
  sizeBits: number;
  hashFunctions: number;
}

const MAXIMUM_HASH_FUNCTIONS = 1000;
const FNV1A_PRIME = BigInt("1099511628211");
const SLASH_UNICODE_VALUE = BigInt("/".charCodeAt(0));
export const PATH_HASH_BASE_HASH = "14695981039346656037";

const countLeadingZeros = (byte: number) => {
  if (byte > 0xff) {
    throw new Error("Input is larger than a byte");
  }

  // This function is equivalent to the function `LeadingZeros8`:
  // Return value equal to 8 minus the minimum number of bits
  // required to represent `byte`, or 8 if `byte` == 0.
  // https://pkg.go.dev/math/bits#LeadingZeros8
  return byte === 0 ? 8 : 8 - byte.toString(2).length;
};

const getNextHash = (hash: bigint) => {
  return BigInt.asUintN(64, (hash ^ SLASH_UNICODE_VALUE) * FNV1A_PRIME);
};

export const readBloomFilter = (
  fsacProfile: FileSystemAccessProfile,
): BloomFilterReader => {
  const bloomFilter = fsacProfile.bloomFilter;
  let hashFunctions = fsacProfile.bloomFilterHashFunctions;

  const lastByte = bloomFilter.at(bloomFilter.length - 1);
  if (lastByte === undefined) {
    throw new Error("Bloom filter is empty");
  }

  const leadingZeros: number = countLeadingZeros(lastByte);
  if (leadingZeros > 7) {
    throw new Error("Bloom filter's trailing byte is not properly padded");
  }

  const sizeBits: number = bloomFilter.length * 8 - leadingZeros - 1;
  if (sizeBits === 0) {
    throw new Error("Bloom filter has zero bits");
  }

  if (hashFunctions > MAXIMUM_HASH_FUNCTIONS) {
    hashFunctions = MAXIMUM_HASH_FUNCTIONS;
  }

  return {
    bloomFilter: bloomFilter,
    sizeBits: sizeBits,
    hashFunctions: hashFunctions,
  };
};

export const containsPathHashes = (
  r: BloomFilterReader,
  pathHashes: PathHashes,
): boolean => {
  let iterHash = pathHashes.baseHash;
  for (let i = 0; i < r.hashFunctions; ++i) {
    const bit = Number(iterHash % BigInt(r.sizeBits));

    if ((r.bloomFilter[Math.floor(bit >> 3)] & (1 << (bit % 8))) === 0) {
      return false;
    }
    iterHash = getNextHash(iterHash);
  }
  return true;
};

export const generateFileSystemReferenceQueryParams = (
  fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined,
  pathHashes?: PathHashes,
): ParsedUrlQueryInput | undefined => {
  if (fileSystemAccessProfileReference === undefined) {
    return undefined;
  }

  let newPathHash = pathHashes;

  if (newPathHash === undefined) {
    newPathHash = new PathHashes(
      BigInt(fileSystemAccessProfileReference.pathHashesBaseHash),
    );
  }

  return {
    fileSystemAccessProfile: JSON.stringify(
      FileSystemAccessProfileReference.toJSON({
        digest: fileSystemAccessProfileReference.digest,
        pathHashesBaseHash: newPathHash.baseHash.toString(),
      }),
    ),
  };
};

export class PathHashes {
  baseHash: bigint;
  constructor(baseHash?: bigint) {
    this.baseHash = baseHash ? baseHash : BigInt(PATH_HASH_BASE_HASH);
  }

  appendComponent(name: string): PathHashes {
    let hash = (this.baseHash ^ SLASH_UNICODE_VALUE) * FNV1A_PRIME;
    for (const c of name) {
      hash = (hash ^ BigInt(c.charCodeAt(0))) * FNV1A_PRIME;
    }
    return new PathHashes(BigInt.asUintN(64, hash));
  }
}
