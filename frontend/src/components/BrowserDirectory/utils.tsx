import { linkOptions } from "@tanstack/react-router";
import type { GlobalToken } from "antd";
import type { DigestFunction_Value } from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import themeStyles from "@/theme/theme.module.css";
import { BrowserPageType } from "@/types/BrowserPageType";
import type { FileSystemAccessProfileReference } from "@/types/FileSystemAccessProfileReference";
import {
  type BloomFilterReader,
  containsPathHashes,
  PathHashes,
} from "@/utils/bloomFilter";
import { generateBrowserSplat, generateFileUrl } from "@/utils/urlGenerator";
import type { CompareState, DirectoryData, DirectoryRenderData } from "./types";

const calculateCompareState = <T,>(
  baseData?: T,
  compareData?: T,
): CompareState => {
  if (!baseData && !compareData) {
    return "equal"; // Should not happen
  }
  if (!compareData) {
    return "unique";
  }
  if (!baseData) {
    return "missing";
  }
  if (typeof baseData === "object" && typeof compareData === "object") {
    if (JSON.stringify(baseData) === JSON.stringify(compareData)) {
      return "equal";
    }
    return "different";
  }
  return baseData === compareData ? "equal" : "different";
};

export const directoryDataToRenderData = (
  instanceName: string,
  digestFunction: DigestFunction_Value,
  input?: DirectoryData,
  compare?: DirectoryData,
  bloomFilterReader?: BloomFilterReader,
  fileSystemAccessProfileReference?: FileSystemAccessProfileReference,
  path?: string,
  useBloomFilter?: boolean,
): DirectoryRenderData => {
  // TODO: Consider a way to cache hashes intead of recalculating them from scratch each time.
  // When doing so, keep in mind that directories in compare action view with the same hash can have different bloom filters.
  const pathHashes =
    useBloomFilter && bloomFilterReader && path
      ? new PathHashes(
          fileSystemAccessProfileReference?.pathHashesBaseHash
            ? BigInt(fileSystemAccessProfileReference.pathHashesBaseHash)
            : undefined,
        ).appendComponents(path.split("/").slice(1))
      : undefined;

  const calcWillBePrefetched = (currentPathHashes: PathHashes | undefined) => {
    if (
      !useBloomFilter ||
      bloomFilterReader === undefined ||
      currentPathHashes === undefined
    ) {
      return undefined;
    }
    return containsPathHashes(bloomFilterReader, currentPathHashes);
  };
  return {
    directories:
      getSortedNames(input?.directories, compare?.directories).map(
        (directoryName) => {
          const directory = input?.directories.find(
            (d) => d.name === directoryName,
          );
          const compareDirectory = compare?.directories.find(
            (d) => d.name === directoryName,
          );
          return {
            nodeType: "directory",
            name: directoryName,
            compareState: calculateCompareState(
              directory?.digest?.hash,
              compareDirectory?.digest?.hash,
            ),
            sizeBytes: directory?.digest?.sizeBytes,
            compareSizeBytes: compareDirectory?.digest?.sizeBytes,
            permissions: "drwxr-xr-x",
            permissionsCompare: calculateCompareState(
              !!directory,
              !!compareDirectory,
            ),
            hash: directory?.digest?.hash,
            compareHash: compareDirectory?.digest?.hash,
            willBePrefetched: directory
              ? calcWillBePrefetched(pathHashes?.appendComponent(directoryName))
              : undefined,
            linkOptions: directory?.digest
              ? linkOptions({
                  to: "/browser/$",
                  params: {
                    _splat: generateBrowserSplat(
                      instanceName,
                      digestFunction,
                      directory?.digest,
                      BrowserPageType.Directory,
                    ),
                  },
                  search: pathHashes
                    ? {
                        fileSystemAccessProfile: {
                          digest: fileSystemAccessProfileReference?.digest,
                          pathHashesBaseHash: pathHashes
                            .appendComponent(directoryName)
                            .toString(),
                        },
                      }
                    : undefined,
                })
              : undefined,
          };
        },
      ) || [],
    files: getSortedNames(input?.files, compare?.files).map((fileName) => {
      const file = input?.files.find((f) => f.name === fileName);
      const compareFile = compare?.files.find((f) => f.name === fileName);
      const isExecutable = file ? file.isExecutable : compareFile?.isExecutable;
      return {
        name: fileName,
        compareState: calculateCompareState(file?.digest, compareFile?.digest),
        nodeType: "file" as const,
        sizeBytes: file?.digest?.sizeBytes,
        compareSizeBytes: compareFile?.digest?.sizeBytes,
        permissions: `-r-${isExecutable ? "x" : "-"}r-${
          isExecutable ? "x" : "-"
        }r-${isExecutable ? "x" : "-"}`,
        permissionsCompare: calculateCompareState(
          file?.isExecutable,
          compareFile?.isExecutable,
        ),
        linkOptions: file?.digest
          ? generateFileUrl(
              instanceName,
              digestFunction,
              file.digest,
              file.name,
            )
          : undefined,
        willBePrefetched: file
          ? calcWillBePrefetched(pathHashes?.appendComponent(fileName))
          : undefined,
      };
    }),
    symlinks: getSortedNames(input?.symlinks, compare?.symlinks).map(
      (symlinkName) => {
        const symlink = input?.symlinks?.find((s) => s.name === symlinkName);
        const compareSymlink = compare?.symlinks?.find(
          (s) => s.name === symlinkName,
        );
        return {
          name: symlinkName,
          compareState: calculateCompareState(
            symlink?.target,
            compareSymlink?.target,
          ),
          target: symlink?.target || compareSymlink?.target || "",
          nodeType: "symlink",
          targetCompare: calculateCompareState(
            symlink?.target,
            compareSymlink?.target,
          ),
          permissions: "lrwxrwxrwx",
          permissionsCompare: calculateCompareState(
            !!symlink,
            !!compareSymlink,
          ),
          willBePrefetched: symlink
            ? calcWillBePrefetched(pathHashes?.appendComponent(symlinkName))
            : undefined,
        };
      },
    ),
  };
};

export const getSortedNames = <T extends { name: string }>(
  list1?: T[],
  list2?: T[],
): string[] => {
  const mergedSet = new Set<string>();
  list1?.forEach((item) => {
    mergedSet.add(item.name);
  });
  list2?.forEach((item) => {
    mergedSet.add(item.name);
  });
  return Array.from(mergedSet).sort((a, b) => a.localeCompare(b));
};

export const styleMap = (
  state: CompareState | "sub_different" | "diff_with_borders" | "link",
  token: GlobalToken,
): React.CSSProperties => {
  switch (state) {
    case "equal":
      return {
        backgroundColor: token.colorBgLayout,
      };
    case "different":
      return {
        backgroundColor: token.yellow3,
      };
    case "unique":
      return {
        backgroundColor: token.green3,
      };
    case "missing":
      return {
        backgroundColor: token.red3,
      };
    case "sub_different":
      return {
        padding: "4px",
        borderRadius: "8px",
        backgroundColor: token.yellow5,
      };
    case "diff_with_borders":
      return {
        backgroundColor: token.yellow5,
        border: `1px solid ${token.colorWarningBorder}`,
        borderRadius: "5px",
      };
    case "link":
      return {
        color: token.colorLink,
        textDecoration: "underline",
        cursor: "pointer",
      };
    default:
      return {
        backgroundColor: token.colorBgContainer,
      };
  }
};

// Packs an array of booleans into a string
export function boolArrayToString(bools: boolean[]): string {
  if (bools.length === 0) return "";

  // 1 byte can hold 8 booleans
  const bytes = new Uint8Array(Math.ceil(bools.length / 8));

  for (let i = 0; i < bools.length; i++) {
    if (bools[i]) {
      // Use bitwise OR to flip the specific bit to 1
      bytes[Math.floor(i / 8)] |= 1 << (i % 8);
    }
  }

  // Convert bytes to a binary string
  const binStr = Array.from(bytes)
    .map((b) => String.fromCharCode(b))
    .join("");

  // Encode to Base64
  return btoa(binStr);
}

// Unpacks a URL-safe string back into an array of booleans
// Might add extra false bools to the end
export function stringToBoolArray(base64: string): boolean[] {
  if (!base64) return [];

  let binStr: string;
  try {
    binStr = atob(base64);
  } catch (_e) {
    // Fallback if the URL string is corrupted
    console.warn("Failed to decode directory state from URL");
    return [];
  }
  const bools: boolean[] = [];

  for (let i = 0; i < binStr.length; i++) {
    const byte = binStr.charCodeAt(i);
    // Extract exactly 8 bits from each byte
    for (let bit = 0; bit < 8; bit++) {
      // Bitwise AND to check if the bit is 1
      bools.push((byte & (1 << bit)) !== 0);
    }
  }
  return bools;
}

export const formattedFileName = (name: string, willBePrefetched?: boolean) => {
  switch (willBePrefetched) {
    case true:
      return <span className={themeStyles.colorSuccess}>{name}</span>;
    case false:
      return (
        <span className={themeStyles.colorFailure}>
          <s>{name}</s>
        </span>
      );
    case undefined:
      return <span>{name}</span>;
  }
};
