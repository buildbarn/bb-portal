import { Link, type LinkOptions } from "@tanstack/react-router";
import CodeText from "../CodeText";

interface Props {
  text: string;
  link: LinkOptions;
  truncate?: boolean;
}

export const CodeLink: React.FC<Props> = ({ text, link, truncate }) => {
  return (
    <Link {...link}>
      <CodeText>
        {truncate === true && text.length > 11
          ? `${text.slice(0, 8)}[…]`
          : text}
      </CodeText>
    </Link>
  );
};
