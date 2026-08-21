import { AnsiUp } from "ansi_up";
import {
  ANSI_HIGHLIGHT_START,
  ANSI_RESET,
  ansiRegex,
  escapeRegex,
} from "./utils";
export interface Props {
  query: string;
  rowIndex: number;
  line: string;
  matchIndexList: number[];
}
const ansi = new AnsiUp();

function highlightLine(line: string, query: string): string {
  return line
    .replace(ansiRegex(), "")
    .replace(
      escapeRegex(query),
      (match) => `${ANSI_HIGHLIGHT_START}${match}${ANSI_RESET}`,
    );
}

const renderLine = (
  query: string,
  rowIndex: number,
  line: string,
  matchIndexList: number[],
) => {
  const matchSet = new Set(matchIndexList);
  if (!query) return ansi.ansi_to_html(line);
  return matchSet.has(rowIndex)
    ? ansi.ansi_to_html(highlightLine(line, query))
    : ansi.ansi_to_html(line);
};

export const LogRow: React.FC<Props> = ({
  query,
  rowIndex,
  line,
  matchIndexList,
}) => {
  return (
    <span
      // biome-ignore lint/security/noDangerouslySetInnerHtml: Should be reworked
      dangerouslySetInnerHTML={{
        __html: renderLine(query, rowIndex, line, matchIndexList),
      }}
    />
  );
};
