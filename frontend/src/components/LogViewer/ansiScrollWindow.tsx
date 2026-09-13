import type React from "react";
import { useEffect, useRef } from "react";
import { experimental_VGrid as VGrid, type VGridHandle } from "virtua";
import PortalAlert from "../PortalAlert";
import styles from "./index.module.css";
import { LogRow } from "./logRow";

interface Props {
  log: string[];
  query: string;
  matchIndexList: number[];
  currentMatchIndex: number;
}

// Takes a log in ansi style, formats it to HTML, and displays it in a scrollable window with virtualization
const AnsiScrollingWindow: React.FC<Props> = ({
  log,
  query,
  matchIndexList,
  currentMatchIndex,
}) => {
  const vListRef = useRef<VGridHandle>(null);

  useEffect(() => {
    if (!matchIndexList.length) return;
    const rowIndex = matchIndexList[currentMatchIndex];
    vListRef.current?.scrollToIndex?.(rowIndex);
  }, [matchIndexList, currentMatchIndex]);

  useEffect(() => {
    if (vListRef.current) {
      vListRef.current.scrollToIndex(log.length - 1);
    }
  }, [log]);

  if (!log) {
    return (
      <PortalAlert
        message="There is no log information to display"
        type="warning"
        showIcon
        className={styles.alert}
      />
    );
  }

  const LINE_HEIGHT = 16.66; // 14px base font size * 0.85 font-size * 1.4 line-height
  const PADDING_HEIGHT = 14; // (Vertical padding + border) * 2
  const MAX_VISIBLE_LINES = Math.min(27.3, log.length); // Make the top line only partially visible to convey that the view screen is scrollable
  return (
    <pre>
      <VGrid
        ref={vListRef}
        style={{
          height: MAX_VISIBLE_LINES * LINE_HEIGHT + PADDING_HEIGHT,
        }}
        className={styles.scrollWindow}
        row={log.length}
        col={1}
        cellHeight={LINE_HEIGHT}
      >
        {({ rowIndex }) => (
          <LogRow
            key={rowIndex}
            query={query}
            rowIndex={rowIndex}
            line={log[rowIndex]}
            matchIndexList={matchIndexList}
          ></LogRow>
        )}
      </VGrid>
    </pre>
  );
};

export { AnsiScrollingWindow };
