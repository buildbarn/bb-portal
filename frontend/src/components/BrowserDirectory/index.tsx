import { useNavigate } from "@tanstack/react-router";
import { Flex, Space } from "antd";
import React, { useCallback, useEffect, useRef, useState } from "react";
import { casByteStreamClient } from "@/grpc/casByteStreamClient";
import {
  type Action,
  type Digest,
  type DigestFunction_Value,
  Directory,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";
import type { FileSystemAccessProfile } from "@/lib/grpc-client/buildbarn/fsac/fsac";
import type { FileSystemAccessProfileReference } from "@/types/FileSystemAccessProfileReference";
import { type BloomFilterReader, readBloomFilter } from "@/utils/bloomFilter";
import { fetchCasObjectAndParse } from "@/utils/fetchCasObject";
import { BrowserDirectoryButtons } from "./browserDirectoryButtons";
import { TopDirectoryComponent } from "./directoryComponent";
import styles from "./index.module.css";
import type {
  CompareContextProps,
  DirectoryData,
  DirectoryRenderData,
  DirectorySizeData,
  QueueItem,
} from "./types";
import {
  boolArrayToString,
  directoryDataToRenderData,
  getSortedNames,
  stringToBoolArray,
} from "./utils";

const TOP_LEVEL_DIR_NAME = "Top level";

export interface CParams {
  instanceName: string;
  digestFunction: DigestFunction_Value;
  digest: Digest;
  fileSystemAccessProfile: FileSystemAccessProfile | undefined;
  fileSystemAccessProfileReference:
    | FileSystemAccessProfileReference
    | undefined;
  action?: Action;
  reducedActionDigest?: Digest | undefined;
}

interface Params {
  baseData: CParams;
  compareData?: CParams;
  mergedMode?: boolean;
  openDirsString: string | undefined;
  useBloomFilter: boolean;
}

export const BrowserDirectory: React.FC<Params> = ({
  baseData,
  compareData,
  mergedMode,
  openDirsString,
  useBloomFilter,
}) => {
  const [openDirs, setOpenDirs] = useState<Set<string>>(new Set([]));
  const directoryData = useRef<Map<string, Promise<DirectoryData>>>(
    new Map([]),
  );
  const [directorySizeData, setDirectorySizeData] = useState<
    Map<string, DirectorySizeData>
  >(new Map([]));
  const parentHashes = useRef<Map<string, string[]>>(new Map([]));
  const [abortController, setAbortController] = useState<AbortController>(
    new AbortController(),
  );
  const [loading, setLoading] = useState<boolean>(false);

  const openDir = useCallback(async (paths: string[] | string) => {
    if (typeof paths === "string") {
      paths = [paths];
    }
    setOpenDirs((old) => new Set([...old, ...paths]));
  }, []);
  let baseBloomFilterReader: BloomFilterReader | undefined;
  if (baseData.fileSystemAccessProfile) {
    baseBloomFilterReader = readBloomFilter(baseData.fileSystemAccessProfile);
  }
  let compareBloomFilterReader: BloomFilterReader | undefined;
  if (compareData?.fileSystemAccessProfile) {
    compareBloomFilterReader = readBloomFilter(
      compareData.fileSystemAccessProfile,
    );
  }

  const navigate = useNavigate();
  const closeDir = useCallback(
    (dir: string) => {
      abortController.abort();
      setOpenDirs((old) => {
        const next = new Set(old);
        next.delete(dir);
        return next;
      });
    },
    [abortController],
  );

  const recalcSize = useCallback(
    (size: number, hash: string, childDigests: Digest[]) => {
      setDirectorySizeData((old) => {
        const newData = new Map(old);
        const parentQueue: {
          hash: string;
        }[] = [];
        const entry: DirectorySizeData | undefined = newData.get(hash);
        if (!entry) {
          newData.set(hash, {
            size: size,
            totalSize:
              size +
              childDigests.reduce(
                (acc, digest) =>
                  acc +
                  (newData.get(digest.hash)?.size ?? Number(digest.sizeBytes)),
                0,
              ),
            fullyResolved: childDigests.length === 0,
            children: childDigests.map((childDigest) => ({
              hash: childDigest.hash,
              sizeBytes: Number(childDigest.sizeBytes),
            })),
          });
        } else {
          entry.totalSize =
            size +
            entry.children.reduce(
              (acc, child) =>
                acc + (newData.get(child.hash)?.size ?? child.sizeBytes),
              0,
            );
        }
        // a directory hash can have two parent hashes if its hash is the same in base and compare, but the parent directory have different hashes
        parentHashes.current.get(hash)?.forEach((parentHash) => {
          parentQueue.push({
            hash: parentHash,
          });
        });
        let remainingLoops = 1000;
        while (parentQueue.length > 0 && remainingLoops > 0) {
          remainingLoops--;
          const currentQueueEntry = parentQueue.pop();
          if (!currentQueueEntry) {
            console.error("Unexpected empty queue entry");
            break;
          }
          const entry = newData.get(currentQueueEntry.hash);
          if (!entry) {
            console.error("Unexpected missing entry", currentQueueEntry.hash);
            continue;
          }
          let totalSize = entry.size;
          let fullyResolved = true;
          entry.children.forEach((child) => {
            const childEntry = newData.get(child.hash);
            totalSize += childEntry?.totalSize ?? child.sizeBytes;
            if (!childEntry || !childEntry.fullyResolved) {
              fullyResolved = false;
            }
          });
          entry.totalSize = totalSize;
          entry.fullyResolved = fullyResolved;
          parentHashes.current
            .get(currentQueueEntry.hash)
            ?.forEach((parentHash) => {
              parentQueue.push({
                hash: parentHash,
              });
            });
        }
        if (remainingLoops <= 0) {
          console.error(
            "Unexpected infinite loop in directory size calculation",
            parentQueue,
          );
        }
        return newData;
      });
    },
    [],
  );
  const getDirFromDigest = useCallback(
    async (
      digest: Digest,
      instanceName: string,
      digestFunction: DigestFunction_Value,
    ): Promise<DirectoryData> => {
      const existingEntry = directoryData.current.get(digest.hash);
      if (existingEntry) {
        return existingEntry;
      } else {
        const newEntry: Promise<DirectoryData> = fetchCasObjectAndParse(
          casByteStreamClient,
          instanceName,
          digestFunction,
          digest,
          Directory,
        ).then((dir) => ({
          ...dir,
          hash: digest.hash,
        }));
        // Update directoryData with the new entry
        directoryData.current.set(digest.hash, newEntry);
        const awaitedEntry = await newEntry;
        // Update parentHashes with the new entry's child directories
        for (const childDir of awaitedEntry.directories) {
          if (!childDir.digest) {
            continue;
          }
          const parentHash = parentHashes.current.get(childDir.digest.hash);
          if (!parentHash) {
            parentHashes.current.set(childDir.digest.hash, [digest.hash]);
          } else if (!parentHash.includes(digest.hash)) {
            parentHash.push(digest.hash);
          }
        }
        // Update total size calculation
        recalcSize(
          Number(digest.sizeBytes) +
            awaitedEntry.files.reduce(
              (acc, file) => acc + Number(file.digest?.sizeBytes),
              0,
            ),
          digest.hash,
          awaitedEntry.directories
            .map((dir) => dir.digest)
            .filter((dir) => dir !== undefined),
        );
        return awaitedEntry;
      }
    },
    [recalcSize],
  );
  const getDirFromPath = useCallback(
    async (
      dir: string[],
      isBase: boolean,
    ): Promise<DirectoryData | undefined> => {
      const data = isBase ? baseData : compareData;
      if (!data) {
        return undefined;
      }
      const instanceName = data.instanceName;
      const digestFunction = data.digestFunction;
      let currentDir = await getDirFromDigest(
        data.digest,
        instanceName,
        digestFunction,
      );

      if (dir[0] !== TOP_LEVEL_DIR_NAME) {
        console.warn("Invalid directory path");
      }
      for (let i = 1; i < dir.length; i++) {
        const nextDirName = dir[i];
        const nextDigest = currentDir.directories.find(
          (d) => d.name === nextDirName,
        )?.digest;
        if (!nextDigest) {
          return undefined;
        }
        currentDir = await getDirFromDigest(
          nextDigest,
          instanceName,
          digestFunction,
        );
      }

      return currentDir;
    },
    [baseData, compareData, getDirFromDigest],
  );

  const getDualDirData = useCallback(
    (
      dir: string,
    ): Promise<[DirectoryData | undefined, DirectoryData | undefined]> => {
      return Promise.all([
        getDirFromPath(dir.split("/"), true),
        getDirFromPath(dir.split("/"), false),
      ]);
    },
    [getDirFromPath],
  );

  const getRenderData = useCallback(
    async (
      dir: string,
      isBase: boolean,
    ): Promise<DirectoryRenderData | null> => {
      const [base, compare] = await getDualDirData(dir);
      return isBase
        ? directoryDataToRenderData(
            baseData.instanceName,
            baseData.digestFunction,
            base,
            compareData ? compare : base,
            baseBloomFilterReader,
            baseData?.fileSystemAccessProfileReference,
            dir,
            useBloomFilter,
          )
        : compareData
          ? directoryDataToRenderData(
              compareData.instanceName,
              compareData.digestFunction,
              compare,
              base,
              compareBloomFilterReader,
              compareData?.fileSystemAccessProfileReference,
              dir,
              useBloomFilter,
            )
          : null;
    },
    [
      getDualDirData,
      baseData,
      compareData,
      baseBloomFilterReader,
      compareBloomFilterReader,
      useBloomFilter,
    ],
  );
  const getTotalSizeData = (hash: string): DirectorySizeData | undefined => {
    return directorySizeData.get(hash);
  };
  const getLeftdata = useCallback(
    (dir: string) => {
      return getRenderData(dir, true);
    },
    [getRenderData],
  );
  const getRightData = useCallback(
    (dir: string) => {
      return getRenderData(dir, false);
    },
    [getRenderData],
  );

  // Opening and closing buttons
  const closeAll = () => {
    setOpenDirs(new Set([TOP_LEVEL_DIR_NAME]));
    abortController.abort();
  };
  const collapse = (times: number) => {
    abortController.abort();
    const resultingOpenDirs = new Set(openDirs);
    for (let _ = 0; _ < times; _++) {
      resultingOpenDirs.forEach((dir) => {
        const hasOpenChild = Array.from(resultingOpenDirs).some(
          (child) => child.startsWith(dir) && child !== dir,
        );
        if (!hasOpenChild) {
          resultingOpenDirs.delete(dir);
        }
      });
    }
    resultingOpenDirs.add(TOP_LEVEL_DIR_NAME);
    setOpenDirs(resultingOpenDirs);
  };
  const expand = async (times: number = 1, relevantOnly: boolean = false) => {
    if (loading) {
      return;
    }
    setLoading(true);
    const newOpenDirs: string[] = [];
    const stack: QueueItem[] = [
      {
        path: TOP_LEVEL_DIR_NAME,
        baseDigest: baseData.digest,
        compareDigest: compareData?.digest,
        depth: 0,
      },
    ];
    let activeAbortController = abortController;
    if (abortController.signal.aborted) {
      activeAbortController = new AbortController();
      setAbortController(activeAbortController);
    }

    while (!activeAbortController.signal.aborted) {
      const item = stack.pop();
      if (item === undefined) {
        break;
      }

      newOpenDirs.push(item.path);
      if (item.depth >= times) {
        continue;
      }

      const baseEntry = item.baseDigest
        ? await getDirFromDigest(
            item.baseDigest,
            baseData.instanceName,
            baseData.digestFunction,
          )
        : undefined;
      const compareEntry =
        item.compareDigest && compareData
          ? await getDirFromDigest(
              item.compareDigest,
              compareData.instanceName,
              compareData.digestFunction,
            )
          : undefined;
      const childDirectoryNames = getSortedNames(
        baseEntry?.directories,
        compareEntry?.directories,
      );
      childDirectoryNames.forEach((childName) => {
        const baseDigest = baseEntry?.directories?.find(
          (dir) => dir.name === childName,
        )?.digest;
        const compareDigest = compareEntry?.directories?.find(
          (dir) => dir.name === childName,
        )?.digest;
        if (
          relevantOnly &&
          (baseDigest?.hash === compareDigest?.hash ||
            !baseDigest ||
            !compareDigest)
        ) {
          return;
        }
        const path = `${item.path}/${childName}`;

        // Don't start counting depth until we encounter an unopened directory
        const newDepth =
          openDirs.has(path) && item.depth === 0 ? 0 : item.depth + 1;

        stack.push({
          path,
          baseDigest: baseDigest,
          compareDigest: compareDigest,
          depth: newDepth,
        });
      });
    }
    if (!activeAbortController.signal.aborted) {
      openDir(newOpenDirs);
    }
    setLoading(false);
  };

  // Serializes the currently open directory tree into a boolean array.
  // Turn the array into a string and store it in the URL as a search parameter.
  useEffect(() => {
    const stateToBooleanArray = async (): Promise<void> => {
      const result: boolean[] = [];
      const queue: [Directory | undefined, Directory | undefined, string][] = [
        [...(await getDualDirData(TOP_LEVEL_DIR_NAME)), TOP_LEVEL_DIR_NAME],
      ];
      while (queue.length > 0) {
        const currentData = queue.shift();
        if (currentData === undefined) {
          break;
        }

        // unify and deduplicate names
        const allChildNames = new Set([
          ...(currentData[0]?.directories.map((dir) => dir.name) || []),
          ...(currentData[1]?.directories.map((dir) => dir.name) || []),
        ]);

        // Sorting is mandaory to guarantee a deterministic order.
        const sortedChildNames = Array.from(allChildNames).sort();

        for (const childName of sortedChildNames) {
          const newPath = `${currentData[2]}/${childName}`;
          if (!openDirs.has(newPath)) {
            // If closed, record a 0 and stop traversing this branch.
            result.push(false);
            continue;
          }
          result.push(true);
          const baseChild = currentData[0]?.directories.find(
            (dir) => dir.name === childName,
          );
          const compareChild = currentData[1]?.directories.find(
            (dir) => dir.name === childName,
          );
          queue.push([
            // These promises sould resolve instantly since we should already have the data cashed it it exists in openDirs.
            baseChild?.digest
              ? await getDirFromDigest(
                  baseChild.digest,
                  baseData.instanceName,
                  baseData.digestFunction,
                )
              : undefined,
            compareChild?.digest && compareData
              ? await getDirFromDigest(
                  compareChild.digest,
                  compareData.instanceName,
                  compareData.digestFunction,
                )
              : undefined,
            newPath,
          ]);
        }
      }
      navigate({
        from: "/browser/$",
        to: "/browser/$",
        search: (prev): typeof prev => ({
          ...prev,
          openDirs: boolArrayToString(result),
        }),
        replace: true,
        resetScroll: false,
      });
    };
    stateToBooleanArray();
  }, [
    openDirs,
    baseData,
    compareData,
    getDirFromDigest,
    getDualDirData,
    navigate,
  ]);

  // Load state of open directories from search parameter
  useEffect(() => {
    const booleanArrayToState = async (input: boolean[]): Promise<void> => {
      const queue: [Directory | undefined, Directory | undefined, string][] = [
        [...(await getDualDirData(TOP_LEVEL_DIR_NAME)), TOP_LEVEL_DIR_NAME],
      ];
      const preloadQueue: [
        Directory | undefined,
        Directory | undefined,
        string,
      ][] = []; // Preload but don't open
      while (queue.length > 0 && input.length > 0) {
        const currentData = queue.shift();
        if (currentData === undefined) {
          break;
        }

        // unify and deduplicate names
        const allChildNames = new Set([
          ...(currentData[0]?.directories.map((dir) => dir.name) || []),
          ...(currentData[1]?.directories.map((dir) => dir.name) || []),
        ]);

        // Sorting is mandaory to guarantee a deterministic order.
        const sortedChildNames = Array.from(allChildNames).sort();

        for (const childName of sortedChildNames) {
          const isOpen = input.shift();
          if (!isOpen) {
            preloadQueue.push(currentData);
            continue;
          }
          const newPath = `${currentData[2]}/${childName}`;
          openDir(newPath);
          const baseChild = currentData[0]?.directories.find(
            (dir) => dir.name === childName,
          );
          const compareChild = currentData[1]?.directories.find(
            (dir) => dir.name === childName,
          );

          const [baseFetched, compareFetched] = await Promise.all([
            baseChild?.digest
              ? getDirFromDigest(
                  baseChild.digest,
                  baseData.instanceName,
                  baseData.digestFunction,
                )
              : Promise.resolve(undefined),
            compareChild?.digest && compareData
              ? getDirFromDigest(
                  compareChild.digest,
                  compareData.instanceName,
                  compareData.digestFunction,
                )
              : Promise.resolve(undefined),
          ]);

          queue.push([baseFetched, compareFetched, newPath]);
        }
      }
      // preload some of the not opened directories
      let remainingPreloads = 100;
      while (preloadQueue.length > 0 && remainingPreloads > 0) {
        remainingPreloads--;
        const currentData = preloadQueue.shift();
        if (currentData === undefined) {
          break;
        }
        // unify and deduplicate names. Sorting is not needed
        const allChildNames = new Set([
          ...(currentData[0]?.directories.map((dir) => dir.name) || []),
          ...(currentData[1]?.directories.map((dir) => dir.name) || []),
        ]);
        for (const childName of allChildNames) {
          const newPath = `${currentData[2]}/${childName}`;
          const baseChild = currentData[0]?.directories.find(
            (dir) => dir.name === childName,
          );
          const compareChild = currentData[1]?.directories.find(
            (dir) => dir.name === childName,
          );

          const [baseFetched, compareFetched] = await Promise.all([
            baseChild?.digest
              ? getDirFromDigest(
                  baseChild.digest,
                  baseData.instanceName,
                  baseData.digestFunction,
                )
              : Promise.resolve(undefined),
            compareChild?.digest && compareData
              ? getDirFromDigest(
                  compareChild.digest,
                  compareData.instanceName,
                  compareData.digestFunction,
                )
              : Promise.resolve(undefined),
          ]);

          preloadQueue.push([baseFetched, compareFetched, newPath]);
        }
      }
    };
    if (openDirs.size === 0) {
      openDir(TOP_LEVEL_DIR_NAME);
    }
    if (openDirs.size > 0 || !openDirsString) {
      return; // If we already have a state, don't overwrite it
    }
    const openDirBoolArray: boolean[] = stringToBoolArray(openDirsString);
    booleanArrayToState(openDirBoolArray);
  }, [
    baseData,
    compareData,
    getDirFromDigest,
    openDir,
    getDualDirData,
    openDirs,
    openDirsString,
  ]);

  return (
    <>
      <BrowserDirectoryButtons
        closeAll={closeAll}
        collapse={collapse}
        expand={expand}
        top={true}
        loading={loading}
        compare={!!compareData}
      />
      <Flex>
        <Space
          orientation="vertical"
          size="large"
          className={styles.directoryContainer}
        >
          <CompareContext.Provider
            value={{
              getData: getLeftdata,
              openDir,
              closeDir,
              openDirs,
              getTotalSizeData,
              mergedMode: !!mergedMode,
            }}
          >
            <TopDirectoryComponent name={TOP_LEVEL_DIR_NAME} />
          </CompareContext.Provider>
        </Space>
        {!mergedMode && compareData && (
          <Space
            orientation="vertical"
            size="large"
            className={styles.directoryContainer}
          >
            <CompareContext.Provider
              value={{
                getData: getRightData,
                openDir,
                closeDir,
                openDirs,
                getTotalSizeData,
                mergedMode: !!mergedMode,
              }}
            >
              <TopDirectoryComponent name={TOP_LEVEL_DIR_NAME} />
            </CompareContext.Provider>
          </Space>
        )}
      </Flex>
      <BrowserDirectoryButtons
        closeAll={closeAll}
        collapse={collapse}
        expand={expand}
        top={false}
        loading={loading}
        compare={!!compareData}
      />
    </>
  );
};

export const CompareContext = React.createContext<CompareContextProps>({
  closeDir: (_: string) => {},
  openDir: (_: string) => {},
  getData: (_: string) => {
    throw new Error("not implemented");
  },
  openDirs: new Set<string>(),
  getTotalSizeData: (_: string) => undefined,
  mergedMode: false,
});
