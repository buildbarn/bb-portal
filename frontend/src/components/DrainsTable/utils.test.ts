import { expect, test } from "vitest";
import {
  drainKey,
  modifyDrainsPermissionUrl,
  patternToPropertyList,
  workerIdPatternFromFields,
} from "./utils";

test("modifyDrainsPermissionUrl", () => {
  expect(modifyDrainsPermissionUrl("")).toBe(
    "/api/v1/checkPermissions/modifyDrains/",
  );
  expect(modifyDrainsPermissionUrl("foo")).toBe(
    "/api/v1/checkPermissions/modifyDrains/foo",
  );
  expect(modifyDrainsPermissionUrl("foo/bar/baz")).toBe(
    "/api/v1/checkPermissions/modifyDrains/foo/bar/baz",
  );
  expect(modifyDrainsPermissionUrl("my instance/a?b#c")).toBe(
    "/api/v1/checkPermissions/modifyDrains/my%20instance/a%3Fb%23c",
  );
});

test("workerIdPatternFromFields", () => {
  expect(workerIdPatternFromFields([])).toEqual({});
  expect(
    workerIdPatternFromFields([
      { key: "hostname", value: "server123.example.com" },
      { key: "thread", value: "3" },
    ]),
  ).toEqual({ hostname: "server123.example.com", thread: "3" });
});

test("patternToPropertyList", () => {
  expect(patternToPropertyList({})).toEqual([]);
  expect(
    patternToPropertyList({ thread: "3", hostname: "a", datacenter: "x" }),
  ).toEqual([
    { name: "datacenter", value: "x" },
    { name: "hostname", value: "a" },
    { name: "thread", value: "3" },
  ]);
});

test("drainKey", () => {
  expect(drainKey({ thread: "3", hostname: "a" })).toBe(
    drainKey({ hostname: "a", thread: "3" }),
  );
  expect(drainKey({ hostname: "a" })).not.toBe(drainKey({ hostname: "b" }));
});
