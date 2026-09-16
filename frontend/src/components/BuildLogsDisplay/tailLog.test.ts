import { describe, expect, it, vi } from "vitest";
import { fetchTailedLog } from "./tailLog";

const linesOf = (count: number, prefix = "line"): string =>
  Array.from({ length: count }, (_, i) => `${prefix} ${i}\n`).join("");

// Serves `log` the way the log endpoint does, from the requested line onwards.
const serve = (log: string) =>
  vi.fn((startLine: number) =>
    Promise.resolve(log.split("\n").slice(startLine).join("\n")),
  );

describe("fetchTailedLog", () => {
  it("reads the whole log when there is nothing cached", async () => {
    const log = linesOf(500);
    const fetchFrom = serve(log);

    expect(await fetchTailedLog(undefined, fetchFrom)).toBe(log);
    expect(fetchFrom).toHaveBeenCalledExactlyOnceWith(0);
  });

  it("reads the whole log when the cached log is shorter than the re-read", async () => {
    const log = linesOf(100);
    const fetchFrom = serve(log);

    expect(await fetchTailedLog(linesOf(50), fetchFrom)).toBe(log);
    expect(fetchFrom).toHaveBeenCalledExactlyOnceWith(0);
  });

  it("re-reads only the tail of a log that has grown", async () => {
    const log = linesOf(700);
    const fetchFrom = serve(log);

    expect(await fetchTailedLog(linesOf(500), fetchFrom)).toBe(log);
    expect(fetchFrom).toHaveBeenCalledExactlyOnceWith(244);
  });

  it("replaces lines that were erased within the re-read", async () => {
    const cached = linesOf(500);
    // The last 10 lines were a progress block that has since been rewritten.
    const log = linesOf(490) + linesOf(3, "progress");
    const fetchFrom = serve(log);

    expect(await fetchTailedLog(cached, fetchFrom)).toBe(log);
    expect(fetchFrom).toHaveBeenCalledExactlyOnceWith(244);
  });

  it("re-reads the whole log when erased lines reach past the re-read", async () => {
    const cached = linesOf(500);
    // Everything from line 100 onwards is gone, so the cached prefix that the
    // tail would be spliced onto no longer matches.
    const log = linesOf(100) + linesOf(5, "progress");
    const fetchFrom = serve(log);

    expect(await fetchTailedLog(cached, fetchFrom)).toBe(log);
    expect(fetchFrom).toHaveBeenNthCalledWith(1, 244);
    expect(fetchFrom).toHaveBeenNthCalledWith(2, 0);
  });

  it("keeps a partial trailing line out of the cached prefix", async () => {
    const log = `${linesOf(400)}incomplete line without a newline`;
    const fetchFrom = serve(log);

    expect(await fetchTailedLog(linesOf(400), fetchFrom)).toBe(log);
    expect(fetchFrom).toHaveBeenCalledExactlyOnceWith(144);
  });
});
