import { useEffect, useState } from 'react';

interface Size {
  width: number;
  height: number;
}

const useScreenSize = () => {
  const [screenSize, setScreenSize] = useState<Size>(getScreenSize());
  function getScreenSize(): Size {
    return {
      width: window.innerWidth,
      height: window.innerHeight,
    };
  }
  useEffect(() => {
    const updateScreenSize = () => {
      setScreenSize(getScreenSize());
    };
    window.addEventListener('resize', updateScreenSize);
    return () => {
      window.removeEventListener('resize', updateScreenSize);
    };
  }, [screenSize]);
  return screenSize;
};

export default useScreenSize;
