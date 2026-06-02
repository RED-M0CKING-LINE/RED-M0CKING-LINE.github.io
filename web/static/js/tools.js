// Loads tool.wasm and connects the CIDR calculator form
(function () {
  'use strict';
  var statusEl = document.getElementById('cidr-status');
  var form = document.getElementById('cidr-form');
  var input = document.getElementById('cidr-input');
  var result = document.getElementById('cidr-result');

  if (typeof Go !== 'function') {
    statusEl.textContent = 'wasm_exec.js failed to load.';
    statusEl.classList.add('err');
    return;
  }

  var go = new Go();
  fetch('/static/wasm/tool.wasm')
    .then(function (r) {
      if (!r.ok) throw new Error('fetch tool.wasm: ' + r.status);
      return WebAssembly.instantiateStreaming(r, go.importObject);
    })
    .then(function (mod) {
      go.run(mod.instance);
      statusEl.textContent = 'Ready. Enter a CIDR.';
    })
    .catch(function (err) {
      statusEl.textContent = 'WASM load failed: ' + err.message;
      statusEl.classList.add('err');
    });

  form.addEventListener('submit', function (e) {
    e.preventDefault();
    var v = input.value.trim();
    if (!v) return;
    var out;
    try {
      // cidrCalc is exposed by the WASM module via js.Global().Set
      out = window.cidrCalc(v);
    } catch (err) {
      out = { error: String(err) };
    }
    renderResult(out);
  });

  function renderResult(out) {
    result.innerHTML = '';
    if (!out) return;
    if (out.error) {
      statusEl.textContent = out.error;
      statusEl.classList.add('err');
      result.hidden = true;
      return;
    }
    statusEl.textContent = 'OK';
    statusEl.classList.remove('err');
    var rows = [
      ['Network',        out.network],
      ['Broadcast',      out.broadcast],
      ['Netmask',        out.netmask],
      ['Wildcard',       out.wildcard],
      ['First host',     out.firstHost],
      ['Last host',      out.lastHost],
      ['Prefix length',  out.prefix],
      ['Total addresses', out.totalAddrs],
      ['Usable hosts',   out.usableHosts]
    ];
    rows.forEach(function (kv) {
      var dt = document.createElement('dt');
      dt.textContent = kv[0];
      var dd = document.createElement('dd');
      dd.textContent = String(kv[1]);
      result.appendChild(dt);
      result.appendChild(dd);
    });
    result.hidden = false;
  }
})();
