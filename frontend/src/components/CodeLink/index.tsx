import { Link, type LinkOptions } from "@tanstack/react-router";

interface Props {
  text: string;
  link: LinkOptions
}

export const CodeLink: React.FC<Props> = ({ text, link }) => {
  return (
    <Link {...link}>
      <code>
        {text.length > 11
          ? `${text.slice(0, 8)}...`
          : text}
      </code>
    </Link>
  );
};
