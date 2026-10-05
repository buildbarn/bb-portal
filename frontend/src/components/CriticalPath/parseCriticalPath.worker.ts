/// <reference lib="webworker" />

import { type CriticalPath, extractCriticalPath } from "./parseCriticalPath";

export interface ParseCriticalPathRequest {
  id: number;
  url: string;
  isGzip: boolean;
}

export interface ParseCriticalPathResponse {
  id: number;
  ok: boolean;
  events?: CriticalPath;
  error?: string;
  // True when the profile blob was reported missing (404/410). Only this
  // failure is permanent; everything else is transient and retryable.
  notFound?: boolean;
}

self.onmessage = async (event: MessageEvent<ParseCriticalPathRequest>) => {
  const { id, url, isGzip } = event.data;
  try {
    const res = await fetch(url);
    if (!res.ok) {
      throw Object.assign(
        new Error(`Failed to fetch profile: ${res.statusText}`),
        { notFound: res.status === 404 || res.status === 410 },
      );
    }
    let trace: unknown;
    if (isGzip) {
      const decompressedStream = res.body?.pipeThrough(
        new DecompressionStream("gzip"),
      );
      if (!decompressedStream) {
        throw new Error("Profile response has no body");
      }
      trace = await new Response(decompressedStream).json();
    } else {
      trace = await res.json();
    }
    const events = extractCriticalPath(trace);
    const response: ParseCriticalPathResponse = { id, ok: true, events };
    self.postMessage(response);
  } catch (error) {
    const response: ParseCriticalPathResponse = {
      id,
      ok: false,
      error: error instanceof Error ? error.message : String(error),
      notFound: (error as { notFound?: boolean })?.notFound === true,
    };
    self.postMessage(response);
  }
};
