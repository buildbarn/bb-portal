import { Input } from "antd";
import { useCallback, useEffect, useMemo } from "react";
import { findMatchIndices } from "./utils";

export interface SearchBarProps {
  query: string;
  setQuery: React.Dispatch<React.SetStateAction<string>>;
  matchIndexList: number[];
  setMatchIndexList: React.Dispatch<React.SetStateAction<number[]>>;
  currentMatchIndex: number;
  setCurrentMatchIndex: React.Dispatch<React.SetStateAction<number>>;
  items: string[];
}

export const SearchBar = ({
  query,
  setQuery,
  matchIndexList,
  setMatchIndexList,
  currentMatchIndex,
  setCurrentMatchIndex,
  items,
}: SearchBarProps) => {
  const matchLimit = 10000;
  const foundMatches = useMemo(
    () => findMatchIndices(items, query, matchLimit),
    [query, items],
  );

  const nextMatch = useCallback(() => {
    setCurrentMatchIndex((i) =>
      matchIndexList.length ? (i + 1) % matchIndexList.length : 0,
    );
  }, [setCurrentMatchIndex, matchIndexList.length]);

  const prevMatch = useCallback(() => {
    setCurrentMatchIndex((i) => {
      if (!matchIndexList.length) return -1;

      return (i - 1 + matchIndexList.length) % matchIndexList.length;
    });
  }, [matchIndexList.length, setCurrentMatchIndex]);

  // Prevent default behavior of F3 outisde searchfield focus
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "F3" && event.shiftKey) {
        event.preventDefault();
        prevMatch();
      } else if (event.key === "F3") {
        event.preventDefault();
        nextMatch();
      }
    };

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [nextMatch, prevMatch]);

  useEffect(() => {
    setMatchIndexList(foundMatches);
  }, [foundMatches, setMatchIndexList]);

  return (
    <Input.Search
      placeholder="Search"
      onSearch={() => nextMatch()}
      onChange={(e) => {
        setQuery(e.target.value);
        // Only a new query starts over from the first match; a log that grew
        // while being tailed keeps the match the reader is on.
        setCurrentMatchIndex(0);
      }}
      value={query}
      suffix={
        <span
          style={{
            display: "inline-block",
            minWidth: "45px",
            textAlign: "right",
            fontSize: "12px",
          }}
        >
          {query.length > 0
            ? foundMatches.length < matchLimit
              ? `${matchIndexList.length > 0 ? currentMatchIndex + 1 : 0} / ${matchIndexList.length}`
              : `${currentMatchIndex + 1} / ${matchLimit}+`
            : null}
        </span>
      }
    />
  );
};
