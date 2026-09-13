import { expect, test } from "vitest";
import { flattenUserInfo } from "./flattenUserInfo";

test("flattenUserInfo", () => {
  expect(flattenUserInfo({})).toEqual({});

  expect(
    flattenUserInfo({
      foo: "bar",
      baz: 1,
      nested: {
        object: "hello",
      },
    }),
  ).toEqual({
    foo: "bar",
    baz: "1",
    "nested/object": "hello",
  });
});
