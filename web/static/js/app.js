// Theme toggle via cookie
// Mobile nav toggle
(function () {
  'use strict';

  function setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    // SameSite=Lax so its sent on top-level nav
    document.cookie = 'theme=' + encodeURIComponent(theme) + '; path=/; max-age=31536000; SameSite=Lax';
  }

  function currentTheme() {
    return document.documentElement.getAttribute('data-theme') || 'dark';
  }

  function getTheme() {
    const match = document.cookie.match(/(?:^|; )theme=([^;]*)/);
    return match ? decodeURIComponent(match[1]) : null;
  }

  // Set theme when page is loaded
  document.addEventListener('DOMContentLoaded', function () {
    const storedTheme = getTheme();
    if (storedTheme) {
      setTheme(storedTheme);
    } else {
      setTheme(currentTheme());
    }
  });

  // Toggle theme
  document.addEventListener('click', function (e) {
    var t = e.target.closest('[data-theme-toggle]');
    if (t) {
      setTheme(currentTheme() === 'dark' ? 'light' : 'dark');
      return;
    }
    var nt = e.target.closest('.nav-toggle');
    if (nt) {
      var nav = document.getElementById('primary-nav');
      if (!nav) return;
      var open = nav.classList.toggle('open');
      nt.setAttribute('aria-expanded', open ? 'true' : 'false');
    }
  });

  // Re-apply theme when the system theme changes AND no explicit choice is set
  var mq = window.matchMedia('(prefers-color-scheme: dark)');
  if (mq && mq.addEventListener) {
    mq.addEventListener('change', function (ev) {
      if (!document.cookie.match(/(?:^|; )theme=/)) {
        setTheme(ev.matches ? 'dark' : 'light');
      }
    });
  }
})();
