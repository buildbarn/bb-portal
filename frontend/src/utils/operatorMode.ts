import { useSyncExternalStore } from "react";

export const OPERATOR_MODE_KEY = "buildbarn_operator_mode";

const listeners = new Set<() => void>();

export const getOperatorMode = (): boolean =>
  localStorage.getItem(OPERATOR_MODE_KEY) === "true";

export const setOperatorMode = (value: boolean) => {
  localStorage.setItem(OPERATOR_MODE_KEY, String(value));
  for (const listener of listeners) {
    listener();
  }
};

const subscribe = (listener: () => void) => {
  listeners.add(listener);
  window.addEventListener("storage", listener);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", listener);
  };
};

export const useOperatorMode = (): [boolean, (value: boolean) => void] => [
  useSyncExternalStore(subscribe, getOperatorMode, () => false),
  setOperatorMode,
];
