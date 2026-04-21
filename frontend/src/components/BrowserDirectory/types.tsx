import type { LinkOptions } from "@tanstack/react-router";
import type {
  Digest,
  Directory,
} from "@/lib/grpc-client/build/bazel/remote/execution/v2/remote_execution";

// What we send to the rendering component
export interface DirectoryRenderData {
  directories: (FileRenderData & {
    nodeType: "directory" | "directory_open" | "aborted";
    target?: never;
    targetCompare?: never;
  })[];
  files: (FileRenderData & {
    nodeType: "file" | "aborted";
    target?: never;
    targetCompare?: never;
    hash?: never;
    compareHash?: never;
  })[];
  symlinks: (FileRenderData & {
    nodeType: "symlink" | "aborted";
    target: string;
    targetCompare: CompareState;
    hash?: never;
    compareHash?: never;
  })[];
}

export interface FileRenderData {
  name: string;
  compareState: CompareState;
  sizeBytes?: string;
  compareSizeBytes?: string;
  target?: string;
  targetCompare?: CompareState;
  nodeType: NodeType;
  permissions: string;
  permissionsCompare: CompareState;
  linkOptions?: LinkOptions | string;
  hash?: string;
  compareHash?: string;
  willBePrefetched?: boolean;
}

export interface DirectorySizeData {
  size: number; // directory object size (tiny) + all files (but not directories) it contains
  totalSize: number;
  fullyResolved: boolean;
  children: {
    sizeBytes: number;
    hash: string;
  }[];
}

export interface DirectoryData extends Directory {
  hash: string;
}

export type NodeType =
  | "file"
  | "directory"
  | "symlink"
  | "directory_open"
  | "aborted";
export type CompareState =
  | "equal"
  | "different"
  | "unique"
  | "missing"
  | "aborted";

export interface CompareContextProps {
  getData: (dir: string) => Promise<DirectoryRenderData | null>;
  closeDir: (dir: string) => void;
  openDir: (dir: string) => void;
  openDirs: Set<string>;
  getTotalSizeData: (hash: string) => DirectorySizeData | undefined;
  mergedMode: boolean;
}

export type QueueItem = {
  path: string;
  baseDigest?: Digest;
  compareDigest?: Digest;
  depth: number;
};
