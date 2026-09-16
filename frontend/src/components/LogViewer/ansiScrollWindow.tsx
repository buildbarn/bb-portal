import type React from "react";
import { useEffect, useRef } from "react";
import { experimental_VGrid as VGrid, type VGridHandle } from "virtua";
import PortalAlert from "../PortalAlert";
import styles from "./index.module.css";
import { LogRow } from "./logRow";

const LINE_HEIGHT = 16.66; // 14px base font size * 0.85 font-size * 1.4 line-height
const PADDING_HEIGHT = 14; // (Vertical padding + border) * 2

// Distance from the end of the log within which it still counts as being
// scrolled to the end. Less than one line, so scrolling away by even a single
// line stops the view from following the log.
const FOLLOW_TAIL_THRESHOLD_PX = LINE_HEIGHT / 2;

const isScrolledToEnd = (grid: VGridHandle): boolean =>
  grid.scrollHeight - grid.scrollTop - grid.viewportHeight <=
  FOLLOW_TAIL_THRESHOLD_PX;

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
  // A log that is still being written to should follow its own end, but not
  // drag a reader who has scrolled up away from what they are looking at.
  const followTail = useRef(true);

  // Keyed on the row rather than the match list, so that a log which grew
  // while being tailed does not scroll back to the match already in view.
  const matchRow: number | undefined = matchIndexList[currentMatchIndex];
  useEffect(() => {
    if (matchRow === undefined) return;
    vListRef.current?.scrollToIndex?.(matchRow);
  }, [matchRow]);

  useEffect(() => {
    if (vListRef.current && followTail.current) {
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

  const MAX_VISIBLE_LINES = Math.min(27.3, log.length); // Make the top line only partially visible to convey that the view screen is scrollable
  return (
    <pre>
      <VGrid
        ref={vListRef}
        style={{
          height: MAX_VISIBLE_LINES * LINE_HEIGHT + PADDING_HEIGHT,
        }}
        className={styles.scrollWindow}
        onScroll={() => {
          if (vListRef.current) {
            followTail.current = isScrolledToEnd(vListRef.current);
          }
        }}
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
