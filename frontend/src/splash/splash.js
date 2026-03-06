(function() {
  var env = window.__env__ || {};
  var titleEl = document.getElementById('preloader-title');
  if (titleEl) {
    var companyPrefix = env.companyName ? env.companyName + ' ' : '';
    titleEl.textContent = 'Welcome to the ' + companyPrefix + 'Buildbarn Portal';
  }
})();
