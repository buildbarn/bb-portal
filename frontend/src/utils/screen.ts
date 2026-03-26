import { useEffect, useState } from 'react';

interface Size {
  width: number;
  height: number;
}

function getScreenSize(): Size {
  return {
    width: window.innerWidth,
    height: window.innerHeight,
  };
}

const useScreenSize = () => {
  const [screenSize, setScreenSize] = useState<Size>(getScreenSize());
  useEffect(() => {
    const updateScreenSize = () => {
      setScreenSize(getScreenSize());
    };
    window.addEventListener('resize', updateScreenSize);
    return () => {
      window.removeEventListener('resize', updateScreenSize);
    };
  }, []);
  return screenSize;
};

export default useScreenSize;
