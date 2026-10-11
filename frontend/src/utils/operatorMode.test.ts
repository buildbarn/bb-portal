import { afterEach, expect, test } from "vitest";
import {
  getOperatorMode,
  OPERATOR_MODE_KEY,
  setOperatorMode,
} from "./operatorMode";

afterEach(() => {
  localStorage.clear();
});

test("defaults to false", () => {
  expect(getOperatorMode()).toBe(false);
});

test("set and read", () => {
  setOperatorMode(true);
  expect(localStorage.getItem(OPERATOR_MODE_KEY)).toBe("true");
  expect(getOperatorMode()).toBe(true);
  setOperatorMode(false);
  expect(localStorage.getItem(OPERATOR_MODE_KEY)).toBe("false");
  expect(getOperatorMode()).toBe(false);
});

test("non-true values read as false", () => {
  for (const value of ["1", "TRUE", "yes", ""]) {
    localStorage.setItem(OPERATOR_MODE_KEY, value);
    expect(getOperatorMode()).toBe(false);
  }
});
