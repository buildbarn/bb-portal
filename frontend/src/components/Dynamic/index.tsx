import React, { useEffect, useState } from 'react';
import FullScreenLoader from '@/components/FullScreenLoader';

// This component waits until client render to render children, effectively bypassing SSR in dev mode
const Dynamic = ({ children }: { children: React.ReactNode }) => {
  const [hasMounted, setHasMounted] = useState(false);

  useEffect(() => {
    setHasMounted(true);
  }, []);

  if (!hasMounted) {
    return <FullScreenLoader />;
  }

  return <>{children}</>;
};

export default Dynamic;
