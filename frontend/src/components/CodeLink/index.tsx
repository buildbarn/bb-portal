import Link from "next/link";

interface Props {
  url: string;
  text: string;
  abbreviate?: boolean;
}

const CodeLink: React.FC<Props> = ({ url, text, abbreviate }) => {
  return (
    <Link href={url}>
      <code>
        {abbreviate === true && text.length > 11
          ? `${text.slice(0, 8)}...`
          : text}
      </code>
    </Link>
  );
};

export default CodeLink;
