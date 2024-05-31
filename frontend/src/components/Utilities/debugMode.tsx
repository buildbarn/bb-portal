export const debugMode = (): true | undefined => {
  const urlParams = new URLSearchParams(window.location.search);
  const debug = urlParams.has('debug');
  const debugLocal = sessionStorage.getItem('debug');

  if (!debug && !debugLocal) {
    return undefined;
  }
  sessionStorage.setItem('debug', 'true');
  return true;
}
