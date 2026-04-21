import { expect, test } from "vitest";
import { boolArrayToString, stringToBoolArray } from "./utils";

const testCases: boolean[][] = [
  [],
  [true],
  [false],
  [true, false, true, true],
  [true, false, true, true, true],
  [true, false, true, true, false],
];

test("boolArrayToStringAndBack", () => {
  testCases.forEach((testCase) => {
    const boolArrayString = boolArrayToString(testCase);
    expect(boolArrayString).toBeTypeOf("string");
    const returnedArray = stringToBoolArray(boolArrayString);
    expect(returnedArray.length).greaterThanOrEqual(testCase.length);
    // The returned array might have extra falses in the end. This is okay.
    expect(returnedArray.slice(0, testCase.length)).toEqual(testCase);
    // Any excess bools should be false.
    expect(returnedArray.slice(testCase.length).every((x) => !x)).toBeTruthy();
  });
});
