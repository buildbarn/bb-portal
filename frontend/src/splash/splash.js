(() => {
  const env = window.__env__ || {};
  const titleEl = document.getElementById("preloader-title");
  if (titleEl) {
    const companyPrefix = env.companyName ? `${env.companyName} ` : "";
    titleEl.textContent = `Welcome to the ${companyPrefix}Buildbarn Portal`;
  }
})();
