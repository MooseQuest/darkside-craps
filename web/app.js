// app.js — UI controller: game flow, passkey auth, history, walkthrough.
import { CrapsGame } from '/engine.js';

const $ = (id) => document.getElementById(id);
const fmt = (n) => (n < 0 ? '-$' : '$') + Math.abs(n).toLocaleString();

// ---- API helpers -----------------------------------------------------------
async function api(path, opts = {}) {
  const res = await fetch(path, {
    method: opts.method || 'GET',
    headers: opts.body ? { 'Content-Type': 'application/json' } : {},
    body: opts.body ? JSON.stringify(opts.body) : undefined,
    credentials: 'same-origin',
  });
  const text = await res.text();
  const data = text ? JSON.parse(text) : {};
  if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`);
  return data;
}

// ---- WebAuthn base64url helpers -------------------------------------------
function b64urlToBuf(s) {
  s = s.replace(/-/g, '+').replace(/_/g, '/');
  const pad = s.length % 4 ? '='.repeat(4 - (s.length % 4)) : '';
  const bin = atob(s + pad);
  const buf = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) buf[i] = bin.charCodeAt(i);
  return buf.buffer;
}
function bufToB64url(buf) {
  const bytes = new Uint8Array(buf);
  let bin = '';
  for (const b of bytes) bin += String.fromCharCode(b);
  return btoa(bin).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

// ---- Global UI state -------------------------------------------------------
let game = null;
let bankroll = 1000;
let auth = { authenticated: false, email: null };
const panels = ['startPanel', 'gamePanel', 'summaryPanel', 'historyPanel'];
function show(id) {
  panels.forEach((p) => ($(p).hidden = p !== id));
  window.scrollTo({ top: 0, behavior: 'smooth' });
}

// ---- Dice rendering (pip layout) ------------------------------------------
const PIP_MAP = {
  1: [4], 2: [0, 8], 3: [0, 4, 8], 4: [0, 2, 6, 8],
  5: [0, 2, 4, 6, 8], 6: [0, 2, 3, 5, 6, 8],
};
function renderDie(el, face) {
  el.replaceChildren();
  const pct = [25, 50, 75]; // 3x3 grid, positioned as % so dice scale to any size
  PIP_MAP[face].forEach((slot) => {
    const pip = document.createElement('span');
    pip.className = 'pip';
    pip.style.left = pct[slot % 3] + '%';
    pip.style.top = pct[Math.floor(slot / 3)] + '%';
    el.appendChild(pip);
  });
}

// ---- Scoreboard / advice updates ------------------------------------------
function paint(r) {
  $('scBankroll').textContent = fmt(r.bankroll);
  const net = r.bankroll - game.initialBankroll;
  const nd = $('scNet');
  nd.textContent = (net >= 0 ? '▲ ' : '▼ ') + fmt(net);
  nd.className = 'score-delta ' + (net >= 0 ? 'pos' : 'neg');
  $('scRolls').textContent = r.roll_count;
  $('scWins').textContent = fmt(r.wins);
  $('scLosses').textContent = fmt(r.losses);
  $('scWagered').textContent = fmt(r.total_wagered);

  const pips = $('scStrikes').querySelectorAll('i');
  pips.forEach((p, i) => p.classList.toggle('on', i < r.strikes));

  const [d1, d2] = r.dice_roll;
  renderDie($('die1'), d1);
  renderDie($('die2'), d2);
  $('rollSum').textContent = `${d1} + ${d2} = ${r.roll_sum}`;

  const chips = $('pointChips');
  if (r.established_points.length) {
    chips.innerHTML = r.established_points.map((p) => `<span class="chip">${p}</span>`).join('');
  } else {
    chips.innerHTML = '<span class="muted">none yet</span>';
  }

  $('statusLine').textContent = r.status;
  $('nextBet').textContent = r.next_bet || '';

  $('summaryBtn').hidden = !r.summary_button;
  $('continueBtn').hidden = !r.continue_button;
  $('endRow').hidden = !(r.summary_button || r.continue_button);
}

// ---- Game flow -------------------------------------------------------------
function startGame() {
  game = new CrapsGame(bankroll);
  const base = bankroll === 1000 ? 25 : 50;
  $('scBankroll').textContent = fmt(bankroll);
  $('scNet').textContent = '';
  ['scRolls', 'scWins', 'scLosses', 'scWagered'].forEach((id) => ($(id).textContent = id === 'scRolls' ? '0' : '$0'));
  $('scStrikes').querySelectorAll('i').forEach((p) => p.classList.remove('on'));
  renderDie($('die1'), 1); renderDie($('die2'), 1);
  $('rollSum').textContent = '—';
  $('pointChips').innerHTML = '<span class="muted">none yet</span>';
  $('statusLine').textContent = `Place your $${base} Don't Pass bet, then roll.`;
  $('nextBet').textContent = '';
  $('summaryBtn').hidden = true;
  $('continueBtn').hidden = true;
  $('endRow').hidden = true;
  show('gamePanel');
  maybeWalkthrough();
}

function doRoll() {
  if (!game) return;
  const betChoice = $('placeBet').checked ? 'bet' : 'no_bet';
  $('die1').classList.add('rolling');
  $('die2').classList.add('rolling');
  const r = game.roll(betChoice);
  setTimeout(() => {
    $('die1').classList.remove('rolling');
    $('die2').classList.remove('rolling');
    paint(r);
  }, 260);
}

async function endGame() {
  const s = game.summary();
  const net = s.net;
  const nd = $('sumNet');
  nd.textContent = (net >= 0 ? '+' : '') + fmt(net);
  nd.className = 'summary-net ' + (net >= 0 ? 'pos' : 'neg');
  $('summaryGrid').innerHTML = [
    ['Final bankroll', fmt(s.final_bankroll)],
    ['Rolls', s.roll_count],
    ['Total wagered', fmt(s.total_wagered)],
    ['Won / Lost', `${fmt(s.wins)} / ${fmt(s.losses)}`],
    ['Points hit', s.points_hit],
    ['Strikes', s.strikes],
  ].map(([k, v]) => `<div class="cell"><b>${v}</b><span>${k}</span></div>`).join('');

  const note = $('savedNote');
  const saveBtn = $('saveSignInBtn');
  note.hidden = true; saveBtn.hidden = true;
  if (auth.authenticated) {
    try {
      await api('/api/sessions', { method: 'POST', body: { summary: s } });
      note.textContent = '✓ Saved to your history';
      note.hidden = false;
    } catch (e) { /* non-fatal */ }
  } else {
    saveBtn.hidden = false;
  }
  show('summaryPanel');
}

// ---- Auth ------------------------------------------------------------------
async function refreshAuth() {
  try { auth = await api('/auth/me'); } catch { auth = { authenticated: false }; }
  $('signInBtn').hidden = auth.authenticated;
  $('whoami').hidden = !auth.authenticated;
  if (auth.authenticated) $('whoamiEmail').textContent = auth.email;
}

function openAuth() { resetAuthSteps(); $('authModal').hidden = false; $('email').focus(); }
function closeAuth() { $('authModal').hidden = true; resetAuthSteps(); }
function authErr(msg) { const e = $('authError'); e.textContent = msg; e.hidden = false; }

// Reset the modal to its first step (email + register/login actions).
function resetAuthSteps() {
  $('authError').hidden = true;
  $('authActions').hidden = false;
  $('codeField').hidden = true;
  $('codeActions').hidden = true;
  $('email').readOnly = false;
  $('code').value = '';
}

// Step 1 of registration: request a verification code. If the server has email
// verification disabled it returns sent=false and we go straight to the passkey.
async function startRegister() {
  const email = $('email').value.trim();
  if (!email) return authErr('Enter your email first.');
  $('authError').hidden = true;
  try {
    const res = await api('/auth/register/send-code', { method: 'POST', body: { email } });
    if (res.sent) {
      $('authActions').hidden = true;
      $('codeField').hidden = false;
      $('codeActions').hidden = false;
      $('codeHint').textContent = `We emailed a 6-digit code to ${email}. Enter it to finish.`;
      $('email').readOnly = true;
      $('code').focus();
    } else {
      await registerPasskey(email, '');
    }
  } catch (e) { authErr(friendly(e)); }
}

// Step 2: create the passkey, sending the code (empty when verification is off).
async function registerPasskey(email, code) {
  try {
    const opts = await api('/auth/register/begin', { method: 'POST', body: { email, code } });
    const pk = opts.publicKey;
    pk.challenge = b64urlToBuf(pk.challenge);
    pk.user.id = b64urlToBuf(pk.user.id);
    (pk.excludeCredentials || []).forEach((c) => (c.id = b64urlToBuf(c.id)));
    const cred = await navigator.credentials.create({ publicKey: pk });
    await api('/auth/register/finish', { method: 'POST', body: {
      id: cred.id, rawId: bufToB64url(cred.rawId), type: cred.type,
      response: {
        attestationObject: bufToB64url(cred.response.attestationObject),
        clientDataJSON: bufToB64url(cred.response.clientDataJSON),
      },
      clientExtensionResults: cred.getClientExtensionResults(),
    } });
    closeAuth(); await refreshAuth();
  } catch (e) {
    // Account already exists: drop back to the email step so "Sign in" is offered.
    if (/already has a passkey/i.test(String(e && e.message || e))) {
      resetAuthSteps();
      authErr('You already have a passkey for this email — tap Sign in.');
      return;
    }
    authErr(friendly(e));
  }
}

async function loginPasskey() {
  const email = $('email').value.trim();
  if (!email) return authErr('Enter your email first.');
  try {
    const opts = await api('/auth/login/begin', { method: 'POST', body: { email } });
    const pk = opts.publicKey;
    pk.challenge = b64urlToBuf(pk.challenge);
    (pk.allowCredentials || []).forEach((c) => (c.id = b64urlToBuf(c.id)));
    const asr = await navigator.credentials.get({ publicKey: pk });
    await api('/auth/login/finish', { method: 'POST', body: {
      id: asr.id, rawId: bufToB64url(asr.rawId), type: asr.type,
      response: {
        authenticatorData: bufToB64url(asr.response.authenticatorData),
        clientDataJSON: bufToB64url(asr.response.clientDataJSON),
        signature: bufToB64url(asr.response.signature),
        userHandle: asr.response.userHandle ? bufToB64url(asr.response.userHandle) : null,
      },
      clientExtensionResults: asr.getClientExtensionResults(),
    } });
    closeAuth(); await refreshAuth();
  } catch (e) { authErr(friendly(e)); }
}

function friendly(e) {
  const m = String(e && e.message || e);
  if (/already has a passkey/i.test(m)) return 'You already have a passkey for this email — tap Sign in.';
  if (/NotAllowed/i.test(m)) return 'No passkey was used. New here? Tap Create account.';
  if (/no passkey found/i.test(m)) return 'No passkey found for that email. Tap Create account.';
  return m;
}

// Add another passkey to the signed-in account (multi-device). No email code —
// the session proves identity; the browser blocks duplicates on the same device.
async function addPasskey() {
  try {
    const opts = await api('/auth/passkey/add/begin', { method: 'POST' });
    const pk = opts.publicKey;
    pk.challenge = b64urlToBuf(pk.challenge);
    pk.user.id = b64urlToBuf(pk.user.id);
    (pk.excludeCredentials || []).forEach((c) => (c.id = b64urlToBuf(c.id)));
    const cred = await navigator.credentials.create({ publicKey: pk });
    await api('/auth/passkey/add/finish', { method: 'POST', body: {
      id: cred.id, rawId: bufToB64url(cred.rawId), type: cred.type,
      response: {
        attestationObject: bufToB64url(cred.response.attestationObject),
        clientDataJSON: bufToB64url(cred.response.clientDataJSON),
      },
      clientExtensionResults: cred.getClientExtensionResults(),
    } });
    alert('Passkey added to this account.');
  } catch (e) {
    const m = String(e && e.message || e);
    if (/InvalidState|already registered|exclude/i.test(m)) { alert('This device already has a passkey for your account.'); return; }
    alert('Could not add passkey: ' + friendly(e));
  }
}

async function signOut() { try { await api('/auth/logout', { method: 'POST' }); } catch {} await refreshAuth(); }

// ---- History ---------------------------------------------------------------
// el builds a DOM node with a text child — never interpolates untrusted data
// into HTML, so server values can't become markup.
function el(tag, cls, text) {
  const n = document.createElement(tag);
  if (cls) n.className = cls;
  if (text != null) n.textContent = String(text);
  return n;
}
const num = (v) => (Number.isFinite(+v) ? +v : 0);

async function openHistory() {
  try {
    const { sessions } = await api('/api/history');
    const agg = await api('/api/summary');
    const netTotal = num(agg.net_total);
    const aggEl = $('historyAgg');
    aggEl.replaceChildren();
    const s1 = el('span'); s1.append('Sessions: ', el('b', null, num(agg.count)));
    const s2 = el('span'); s2.append('Net lifetime: ',
      el('b', netTotal >= 0 ? 'pos' : 'neg', (netTotal >= 0 ? '+' : '') + fmt(netTotal)));
    const s3 = el('span'); s3.append('Rolls: ', el('b', null, num(agg.rolls_total)));
    aggEl.append(s1, s2, s3);

    const list = $('historyList');
    list.replaceChildren();
    if (!sessions.length) {
      list.append(el('p', 'muted', 'No saved sessions yet. Finish a game while signed in.'));
    } else {
      for (const s of sessions) {
        const net = num(s.net);
        const when = s.created_at ? new Date(s.created_at).toLocaleString() : '';
        const left = el('div');
        left.append(el('b', null, fmt(num(s.initial_bankroll))), ' → ' + fmt(num(s.final_bankroll)),
          el('br'), el('small', null, `${num(s.roll_count)} rolls · ${when}`));
        const right = el('div', 'net ' + (net >= 0 ? 'pos' : 'neg'), (net >= 0 ? '+' : '') + fmt(net));
        const row = el('div', 'hrow'); row.append(left, right);
        list.append(row);
      }
    }
    show('historyPanel');
  } catch (e) { alert('Could not load history: ' + e.message); }
}

// ---- Walkthrough (first run) ----------------------------------------------
const WALK = [
  ['Welcome', 'The Long Dark Side', "You're betting the Don't Pass line — with the house, against the shooter. The goal is a long, disciplined session, not a big score."],
  ['Bankroll', 'Base bet scales', 'A $1,000 bankroll bets $25 a round; $2,000 bets $50. Everything — the lay odds too — scales from there.'],
  ['Rolls', 'Roll & follow the advice', 'Hit “Roll dice”. After each roll the status line tells you what happened and the gold line tells you exactly what to lay next.'],
  ['Strikes', 'Three strikes = pause', "Each time a working point repeats you take a strike and lose that lay. After three strikes, stop betting and wait for a 7."],
  ['Seven', 'The 7 pays you', 'Seven-out clears the table in your favor — you win every working lay and the count resets. Discipline over adrenaline.'],
];
let walkI = 0;
function maybeWalkthrough() {
  if (localStorage.getItem('walk_done')) return;
  walkI = 0; renderWalk(); $('walkModal').hidden = false;
}
function renderWalk() {
  const [step, title, body] = WALK[walkI];
  $('walkStep').textContent = `Step ${walkI + 1} of ${WALK.length} · ${step}`;
  $('walkTitle').textContent = title;
  $('walkBody').textContent = body;
  $('walkDots').innerHTML = WALK.map((_, i) => `<i class="${i === walkI ? 'on' : ''}"></i>`).join('');
  $('walkNext').textContent = walkI === WALK.length - 1 ? 'Got it ✓' : 'Next ▸';
}
function walkNext() {
  if (walkI < WALK.length - 1) { walkI++; renderWalk(); }
  else closeWalk();
}
function closeWalk() { localStorage.setItem('walk_done', '1'); $('walkModal').hidden = true; }

// ---- Wiring ----------------------------------------------------------------
function wire() {
  $('bankrollSeg').querySelectorAll('.seg').forEach((b) => b.addEventListener('click', () => {
    $('bankrollSeg').querySelectorAll('.seg').forEach((s) => { s.classList.remove('active'); s.setAttribute('aria-checked', 'false'); });
    b.classList.add('active'); b.setAttribute('aria-checked', 'true');
    bankroll = parseInt(b.dataset.val, 10);
  }));
  $('startBtn').addEventListener('click', startGame);
  $('rollBtn').addEventListener('click', doRoll);
  $('summaryBtn').addEventListener('click', endGame);
  $('continueBtn').addEventListener('click', () => { $('summaryBtn').hidden = true; $('continueBtn').hidden = true; $('endRow').hidden = true; });
  $('playAgainBtn').addEventListener('click', () => show('startPanel'));

  $('signInBtn').addEventListener('click', openAuth);
  $('startSignInLink').addEventListener('click', (e) => { e.preventDefault(); openAuth(); });
  $('saveSignInBtn').addEventListener('click', openAuth);
  $('authClose').addEventListener('click', closeAuth);
  $('passkeyRegisterBtn').addEventListener('click', startRegister);
  $('verifyCreateBtn').addEventListener('click', () => registerPasskey($('email').value.trim(), $('code').value.trim()));
  $('codeBackBtn').addEventListener('click', resetAuthSteps);
  $('passkeyLoginBtn').addEventListener('click', loginPasskey);
  $('signOutBtn').addEventListener('click', signOut);
  $('addPasskeyBtn').addEventListener('click', addPasskey);
  $('historyBtn').addEventListener('click', openHistory);
  $('historyCloseBtn').addEventListener('click', () => show('startPanel'));

  $('walkNext').addEventListener('click', walkNext);
  $('walkSkip').addEventListener('click', closeWalk);

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !$('gamePanel').hidden && $('summaryBtn').hidden) doRoll();
    if (e.key === 'Escape') { closeAuth(); }
  });

  if (!window.PublicKeyCredential) {
    $('signInBtn').title = 'Passkeys not supported in this browser';
  }
}

wire();
refreshAuth();
