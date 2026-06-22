// Inline theme bootstrap: avoid FOUC. Reads system pref and theme override stored in a cookie set by the toggle
(function () {
  try {
    var m = document.cookie.match(/(?:^|; )theme=([^;]+)/);
    var t = m ? decodeURIComponent(m[1]) : null;
    if (!t) {
      t = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }
    document.documentElement.setAttribute('data-theme', t);
  } catch (e) {}
})();