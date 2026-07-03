// engine.test.mjs — run with: node --test
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { CrapsGame } from './engine.js';

// Helper: build a game whose dice come from a scripted queue of [d1,d2] pairs.
function scripted(bankroll, pairs) {
  const q = [...pairs];
  const roller = () => {
    if (!q.length) throw new Error('scripted roller exhausted');
    return q.shift();
  };
  return new CrapsGame(bankroll, roller);
}

test('initial state: Don\'t Pass base bet by bankroll tier', () => {
  assert.equal(new CrapsGame(1000).currentBets[0].bet, 25);
  assert.equal(new CrapsGame(2000).currentBets[0].bet, 50);
});

test('come-out 7 loses the Don\'t Pass bet', () => {
  const g = scripted(1000, [[3, 4]]); // 7
  const r = g.roll('bet');
  assert.equal(r.roll_sum, 7);
  assert.equal(g.bankroll, 975);
  assert.equal(r.losses, 25);
  assert.ok(r.summary_button && r.continue_button);
  assert.equal(g.comeOutRoll, true);
});

test('come-out 11 loses the Don\'t Pass bet', () => {
  const g = scripted(1000, [[5, 6]]);
  const r = g.roll('bet');
  assert.equal(g.bankroll, 975);
  assert.equal(r.roll_sum, 11);
});

test('come-out 2/3 wins the Don\'t Pass bet', () => {
  const g = scripted(1000, [[1, 1]]); // 2
  const r = g.roll('bet');
  assert.equal(g.bankroll, 1025);
  assert.equal(r.wins, 25);
  assert.ok(r.summary_button);
});

test('come-out 12 is a push', () => {
  const g = scripted(1000, [[6, 6]]);
  const r = g.roll('bet');
  assert.equal(g.bankroll, 1000);
  assert.match(r.status, /push/);
});

test('come-out point establishes + places lay odds', () => {
  const g = scripted(1000, [[2, 2]]); // point 4
  const r = g.roll('bet');
  assert.equal(g.comeOutRoll, false);
  assert.deepEqual(g.establishedPoints, [4]);
  const lay = g.currentBets.find((b) => b.type === 'Lay Odds' && b.point === 4);
  assert.ok(lay, 'lay odds placed on 4');
  assert.equal(lay.bet, 90);
  assert.match(r.status, /Point established/);
});

test('seven-out after a point wins all lay bets and resets', () => {
  // point 4 (lay potential_win 45), then 7-out
  const g = scripted(1000, [[2, 2], [3, 4]]);
  g.roll('bet');                 // establish 4
  const r = g.roll('bet');       // 7 out
  assert.equal(r.roll_sum, 7);
  assert.match(r.status, /won on all points/);
  assert.equal(g.bankroll, 1045); // +45 lay win
  assert.equal(g.comeOutRoll, true);
  assert.deepEqual(g.establishedPoints, []);
  assert.equal(g.strikes, 0);
  assert.equal(g.currentBets.length, 1); // reset to Don't Pass
});

test('point hit is a strike and loses the lay bet', () => {
  // establish 4 (lay bet 90), then hit 4 again with no new bet placed
  const g = scripted(1000, [[2, 2], [2, 2]]);
  g.roll('bet');                     // establish 4, bankroll 1000
  const r = g.roll('no_bet');        // hit 4, no fresh Don't Come
  assert.equal(g.strikes, 1);
  assert.equal(g.pointsHit, 1);
  assert.equal(g.bankroll, 1000 - 90); // lost the 90 lay bet only
  assert.match(r.status, /strike 1/);
  assert.ok(!g.establishedPoints.includes(4));
});

test('three strikes pauses betting (bets cleared)', () => {
  // Establish and immediately re-hit three different points.
  // Rolls: est 4, hit 4, est 5, hit 5, est 6, hit 6
  const g = scripted(1000, [
    [2, 2], [2, 2],   // 4 then 4
    [2, 3], [2, 3],   // 5 then 5
    [3, 3], [3, 3],   // 6 then 6
  ]);
  g.roll('bet'); g.roll('no_bet');   // strike 1
  g.roll('bet'); g.roll('no_bet');   // strike 2
  g.roll('bet');
  const r = g.roll('no_bet');        // strike 3
  assert.equal(g.strikes, 3);
  assert.match(r.status, /Three strikes/);
  assert.equal(g.currentBets.length, 0);
});

test('lay guidance suppressed on non-point come-out numbers', () => {
  const g = scripted(1000, [[3, 4]]); // 7
  const r = g.roll('bet');
  assert.equal(r.next_bet, '');
});

test('lay guidance shown for a point number on come-out', () => {
  const g = scripted(1000, [[5, 5]]); // 10
  const r = g.roll('bet');
  assert.match(r.next_bet, /Lay \$90.*4 or 10/);
});

test('summary reports net and totals', () => {
  const g = scripted(1000, [[1, 1]]); // win 25
  g.roll('bet');
  const s = g.summary();
  assert.equal(s.final_bankroll, 1025);
  assert.equal(s.net, 25);
  assert.equal(s.roll_count, 1);
  assert.equal(s.total_wagered, 25);
});

test('$2000 tier uses $50 base and larger lay odds', () => {
  const g = scripted(2000, [[2, 2]]); // point 4
  g.roll('bet');
  const lay = g.currentBets.find((b) => b.type === 'Lay Odds');
  assert.equal(lay.bet, 100);
  assert.equal(lay.potential_win, 100);
});
