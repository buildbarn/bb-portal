import type { FileDetailsFragment } from "@/graphql/__generated__/graphql";
import { digestFunctionValueFromString } from "@/utils/digestFunctionUtils";
import { generateFileUrl } from "@/utils/urlGenerator";
import type { CriticalPath } from "./parseCriticalPath";
import type {
  ParseCriticalPathRequest,
  ParseCriticalPathResponse,
} from "./parseCriticalPath.worker";

const isGzipProfile = (profile: FileDetailsFragment) =>
  profile.filePath.path.endsWith(".gz");

const getProfileUrl = (profile: FileDetailsFragment) =>
  generateFileUrl(
    profile.digest.rev2InstanceName,
    digestFunctionValueFromString(profile.digest.digestFunction),
    {
      hash: profile.digest.hash,
      sizeBytes: profile.digest.sizeBytes.toString(),
    },
    profile.filePath.path,
  );

// Errors reported by a worker for its own request that must not be
// retried: the profile blob is permanently gone (404/410).
class WorkerDataError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "WorkerDataError";
  }
}

// Profile blob missing. Cached as undefined so the chart stays hidden
// without hammering the CAS on remounts.
class ProfileNotFoundError extends Error {}

interface PendingRequest {
  resolve: (events: CriticalPath) => void;
  reject: (error: unknown) => void;
}

interface WorkerEntry {
  worker: Worker;
  pending: Map<number, PendingRequest>;
}

// The parsing pipeline is CPU-bound on the worker thread, so a single
// worker would serialize all requests. Use a small pool, sized by the
// number of logical CPUs, to parse profiles concurrently. The pool is
// capped to keep the transient memory of parallel profile parses in
// check; each parse temporarily holds a decompressed trace and its
// parsed representation.
const MAX_WORKER_POOL_SIZE = 4;

const workerEntries: WorkerEntry[] = [];
let nextDispatchIndex = 0;
let nextRequestId = 0;

const createWorkerEntry = (): WorkerEntry => {
  const worker = new Worker(
    new URL("./parseCriticalPath.worker.ts", import.meta.url),
    { type: "module" },
  );
  const entry: WorkerEntry = { worker, pending: new Map() };
  worker.onmessage = (event: MessageEvent<ParseCriticalPathResponse>) => {
    const { id, ok, events, error, notFound } = event.data;
    const pending = entry.pending.get(id);
    if (!pending) {
      return;
    }
    entry.pending.delete(id);
    if (ok && events) {
      pending.resolve(events);
    } else if (notFound) {
      pending.reject(
        new ProfileNotFoundError(error ?? "Profile blob not found"),
      );
    } else {
      pending.reject(new WorkerDataError(error ?? "Unknown worker error"));
    }
  };
  worker.onerror = (event) => {
    // The worker infrastructure failed (e.g. the script could not load).
    // Reject its requests as infrastructure errors, and drop the worker
    // so that a new one is created for subsequent requests.
    const entryIndex = workerEntries.indexOf(entry);
    if (entryIndex !== -1) {
      workerEntries.splice(entryIndex, 1);
    }
    const error = new Error(event.message || "Worker failed");
    for (const pending of entry.pending.values()) {
      pending.reject(error);
    }
    entry.pending.clear();
  };
  return entry;
};

const getWorkerEntry = (): WorkerEntry => {
  const desiredPoolSize = Math.min(
    MAX_WORKER_POOL_SIZE,
    Math.max(1, navigator.hardwareConcurrency ?? 1),
  );
  if (workerEntries.length < desiredPoolSize) {
    const entry = createWorkerEntry();
    workerEntries.push(entry);
    return entry;
  }
  const entry = workerEntries[nextDispatchIndex % workerEntries.length];
  nextDispatchIndex += 1;
  return entry;
};

export const getCriticalPath = async (
  profile: FileDetailsFragment,
): Promise<CriticalPath | undefined> => {
  const url = getProfileUrl(profile);
  const request: ParseCriticalPathRequest = {
    id: nextRequestId++,
    url,
    isGzip: isGzipProfile(profile),
  };
  const { worker, pending } = getWorkerEntry();
  try {
    return await new Promise<CriticalPath>((resolve, reject) => {
      pending.set(request.id, { resolve, reject });
      worker.postMessage(request);
    });
  } catch (error) {
    if (error instanceof ProfileNotFoundError) {
      // The profile blob is gone permanently; suppress the error so the
      // chart simply stays hidden.
      return undefined;
    }
    // Transient failure (network error, 5xx, parse failure). Rethrow so
    // that react-query applies its retry policy and refetches the query
    // on remount, instead of caching the failure permanently.
    throw error;
  }
};
