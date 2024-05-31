'use client';
import { useState, type PropsWithChildren } from 'react';
import { useServerInsertedHTML } from 'next/navigation';
import { createCache, extractStyle, StyleProvider } from '@ant-design/cssinjs';

// RootStyleRegistry is used to circumvent flash of unstyled content when using AntD with NextJS, Next uses SSR in Dev mode,
// so styles need to be extracted and loaded on initial page load.
const RootStyleRegistry = ({ children }: PropsWithChildren) => {
  const [cache] = useState(() => createCache());

  useServerInsertedHTML(() => {
    return (
      <script
        dangerouslySetInnerHTML={{
          __html: `</script>${extractStyle(cache)}<script>`,
        }}
      />
    );
  });

  return <StyleProvider cache={cache}>{children}</StyleProvider>;
};

export default RootStyleRegistry;
