type Props = {
  url?: string;
  children: React.ReactNode;
};

export const OptionalLinkWrapper: React.FC<Props> = ({ url, children }) => {
  if (url) {
    return <a href={url}>{children}</a>;
  }
  return children;
};
