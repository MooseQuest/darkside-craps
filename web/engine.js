// engine.js — Dark Side Craps strategy engine (client-side, dependency-free).
//
// Faithful port of the Python CrapsStateMachine (app.py). The engine is a pure
// state machine: feed it "roll" events and it returns a response snapshot. Dice
// are produced by an injectable roller so tests can script outcomes.
//
// Usage:
//   import { CrapsGame } from './engine.js'
//   const g = new CrapsGame(1000)
//   const r = g.roll('bet')            // random dice
//   const r = g.roll('bet', [3, 4])    // forced dice (tests)

const POINT_NUMBERS = [4, 5, 6, 8, 9, 10];

function randomDie() {
  // 1..6 inclusive
  return Math.floor(Math.random() * 6) + 1;
}

export class CrapsGame {
  // roller: optional () => [d1, d2] for deterministic tests.
  constructor(initialBankroll, roller = null) {
    this.initialBankroll = initialBankroll;
    this.bankroll = initialBankroll;
    this.totalWagered = 0;
    const base = this._baseBet();
    this.currentBets = [{ type: "Don't Pass", bet: base, potential_win: base }];
    this.establishedPoints = [];
    this.results = [];
    this.rollCount = 0;
    this.pointsHit = 0;
    this.pointsLost = 0;
    this.wins = 0;
    this.losses = 0;
    this.strikes = 0;
    this.comeOutRoll = true;
    this._roller = roller || (() => [randomDie(), randomDie()]);
  }

  _baseBet() {
    // $25 base for a $1000 bankroll, $50 otherwise (mirrors the Python `== 1000`).
    return this.initialBankroll === 1000 ? 25 : 50;
  }

  // Public entry point. event is always 'roll'; betChoice is 'bet' | 'no_bet'.
  roll(betChoice = 'bet', dice = null) {
    return this.comeOutRoll
      ? this._comeOutRoll(betChoice, dice)
      : this._pointEstablishedRoll(betChoice, dice);
  }

  _rollDice(dice) {
    if (dice && dice.length === 2) return [dice[0], dice[1]];
    return this._roller();
  }

  _comeOutRoll(betChoice, dice) {
    const diceRoll = this._rollDice(dice);
    const rollSum = diceRoll[0] + diceRoll[1];
    const initialBet = this._baseBet();
    const [, layMessage] = this._getLayBet(rollSum, true);
    const response = {
      status: 'continue',
      roll_sum: rollSum,
      dice_roll: diceRoll,
      next_bet: layMessage,
    };

    if (betChoice === 'bet') {
      this.totalWagered += initialBet;
      if (rollSum === 7 || rollSum === 11) {
        this.bankroll -= initialBet;
        this.losses += initialBet;
        response.status = "Seven or Eleven rolled, you lost the Don't Pass bet";
        response.summary_button = true;
        response.continue_button = true;
        this.comeOutRoll = true;
        this.currentBets = [{ type: "Don't Pass", bet: initialBet, potential_win: initialBet }];
      } else if (rollSum === 2 || rollSum === 3) {
        this.bankroll += initialBet;
        this.wins += initialBet;
        response.status = "Two or Three rolled, you won the Don't Pass bet";
        response.summary_button = true;
        this.currentBets = [{ type: "Don't Pass", bet: initialBet, potential_win: initialBet }];
      } else if (rollSum === 12) {
        response.status = "Twelve rolled, it's a push";
      } else {
        this.establishedPoints.push(rollSum);
        this.comeOutRoll = false;
        this._placeOddsBet(rollSum);
        response.status = 'Point established. Place a $' + initialBet + " Don't Come bet.";
      }
    } else {
      response.status = 'No bet';
    }

    this.results.push(diceRoll);
    this.rollCount += 1;
    this._snapshot(response);
    return response;
  }

  _pointEstablishedRoll(betChoice, dice) {
    const diceRoll = this._rollDice(dice);
    const rollSum = diceRoll[0] + diceRoll[1];
    const initialBet = this._baseBet();
    const [, layMessage] = this._getLayBet(rollSum, false);
    const response = {
      status: 'continue',
      roll_sum: rollSum,
      dice_roll: diceRoll,
      next_bet: layMessage,
    };
    const placingBet = betChoice === 'bet';
    if (placingBet) this.totalWagered += initialBet;

    if (rollSum === 7) {
      this._resolveSevenOut(initialBet, response);
    } else if (this.establishedPoints.includes(rollSum)) {
      // The fresh Don't Come bet only loses if one was placed this roll.
      const extraLoss = placingBet ? initialBet : 0;
      this._resolvePointHit(rollSum, extraLoss, response);
    } else if (POINT_NUMBERS.includes(rollSum) && placingBet) {
      this.establishedPoints.push(rollSum);
      this._placeOddsBet(rollSum);
      response.status = 'Point established. Place a $' + initialBet + " Don't Come bet.";
    } else {
      response.status = placingBet ? 'No effect on existing bets' : 'No bet';
    }

    this.results.push(diceRoll);
    this.rollCount += 1;
    this._snapshot(response);
    return response;
  }

  _resolveSevenOut(initialBet, response) {
    const totalWins = this.currentBets
      .filter((b) => this.establishedPoints.includes(b.point))
      .reduce((acc, b) => acc + b.potential_win, 0);
    this.bankroll += totalWins;
    this.wins += totalWins;
    response.status = 'Seven rolled, you won on all points';
    response.summary_button = true;
    response.continue_button = true;
    this.comeOutRoll = true;
    this.establishedPoints = [];
    this.currentBets = [{ type: "Don't Pass", bet: initialBet, potential_win: initialBet }];
    this.strikes = 0;
  }

  _resolvePointHit(rollSum, extraLoss, response) {
    const pointBet = this.currentBets.find((b) => b.point === rollSum) || null;
    const loss = extraLoss + (pointBet ? pointBet.bet : 0);
    this.bankroll -= loss;
    this.losses += loss;
    this.pointsHit += 1;
    this.establishedPoints = this.establishedPoints.filter((p) => p !== rollSum);
    if (pointBet) this.currentBets = this.currentBets.filter((b) => b !== pointBet);
    response.status =
      "Point hit. Don't Behind " + rollSum +
      " and it's replaced by your Don't Come bet. This is strike " +
      (this.strikes + 1) + ". Place another Don't Come bet.";
    this.strikes += 1;
    if (this.strikes >= 3) {
      response.status += ' Three strikes, wait for a 7 before betting again.';
      this.currentBets = [];
    }
  }

  _placeOddsBet(rollSum) {
    const k = this.initialBankroll === 1000;
    if (rollSum === 4 || rollSum === 10) {
      this.currentBets.push({ type: 'Lay Odds', point: rollSum, bet: k ? 90 : 100, potential_win: k ? 45 : 100 });
    } else if (rollSum === 5 || rollSum === 9) {
      this.currentBets.push({ type: 'Lay Odds', point: rollSum, bet: k ? 60 : 75, potential_win: k ? 40 : 75 });
    } else if (rollSum === 6 || rollSum === 8) {
      this.currentBets.push({ type: 'Lay Odds', point: rollSum, bet: k ? 30 : 60, potential_win: k ? 25 : 60 });
    }
  }

  _getLayBet(rollSum, onComeOut) {
    if (!POINT_NUMBERS.includes(rollSum)) return [0, ''];
    if (!onComeOut && this.establishedPoints.includes(rollSum)) return [0, ''];
    const k = this.initialBankroll === 1000;
    if (rollSum === 4 || rollSum === 10)
      return [k ? 90 : 100, 'Lay $90 if bankroll is $1000 or $100 if bankroll is $2000 on 4 or 10'];
    if (rollSum === 5 || rollSum === 9)
      return [k ? 60 : 75, 'Lay $60 if bankroll is $1000 or $75 if bankroll is $2000 on 5 or 9'];
    return [k ? 30 : 60, 'Lay $30 if bankroll is $1000 or $60 if bankroll is $2000 on 6 or 8'];
  }

  _snapshot(response) {
    response.current_bets = this.currentBets.map((b) => ({ ...b }));
    response.bankroll = this.bankroll;
    response.results = this.results.map((r) => [...r]);
    response.roll_count = this.rollCount;
    response.wins = this.wins;
    response.losses = this.losses;
    response.established_points = [...this.establishedPoints];
    response.strikes = this.strikes;
    response.total_wagered = this.totalWagered;
    response.come_out_roll = this.comeOutRoll;
  }

  // Serializable summary for persistence / the summary screen.
  summary() {
    return {
      initial_bankroll: this.initialBankroll,
      final_bankroll: this.bankroll,
      roll_count: this.rollCount,
      results: this.results.map((r) => [...r]),
      points_established: this.establishedPoints.length,
      points_hit: this.pointsHit,
      points_lost: this.pointsLost,
      wins: this.wins,
      losses: this.losses,
      strikes: this.strikes,
      total_wagered: this.totalWagered,
      net: this.bankroll - this.initialBankroll,
    };
  }
}

export const _internal = { POINT_NUMBERS };
