import random
from flask import Flask, render_template, request, jsonify, session
from flask_session import Session

app = Flask(__name__)
app.config['SECRET_KEY'] = 'supersecretkey'
app.config['SESSION_TYPE'] = 'filesystem'
Session(app)

# Function to simulate a dice roll
def roll_dice():
    return random.randint(1, 6), random.randint(1, 6)

# Initialize game state
def init_game(initial_bankroll):
    session['initial_bankroll'] = initial_bankroll
    session['bankroll'] = initial_bankroll
    session['current_bets'] = [{'type': 'Don\'t Pass', 'bet': 25 if initial_bankroll == 1000 else 50, 'potential_win': 25 if initial_bankroll == 1000 else 50}]
    session['established_points'] = []
    session['results'] = []
    session['roll_count'] = 0
    session['points_hit'] = 0
    session['wins'] = 0
    session['losses'] = 0
    session['strikes'] = 0

# Function to simulate a single roll in the game
def game_roll():
    roll = roll_dice()
    roll_sum = sum(roll)
    session['results'].append(roll_sum)
    session['roll_count'] += 1
    return roll_sum

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/start', methods=['POST'])
def start_game():
    initial_bankroll = int(request.form['bankroll'])
    init_game(initial_bankroll)
    session['bankroll'] -= 25 if initial_bankroll == 1000 else 50  # Initial bet on Don't Pass
    return jsonify({
        'status': 'Game started',
        'initial_bankroll': initial_bankroll,
        'current_bets': session['current_bets'],
        'bankroll': session['bankroll']
    })

@app.route('/roll', methods=['POST'])
def roll():
    roll_sum = game_roll()
    initial_bet = 25 if session['initial_bankroll'] == 1000 else 50
    response = {
        'status': 'continue',
        'roll_sum': roll_sum,
        'current_bets': session['current_bets'],
        'bankroll': session['bankroll'],
        'results': session['results'],
        'roll_count': session['roll_count'],
        'wins': session['wins'],
        'losses': session['losses']
    }

    if len(session['established_points']) < 4:
        if roll_sum in [7, 11]:
            session['bankroll'] -= initial_bet
            session['losses'] += initial_bet
            response['status'] = 'Seven or Eleven rolled, you lost the Don\'t Pass bet'
        elif roll_sum in [2, 3]:
            session['bankroll'] += initial_bet
            session['wins'] += initial_bet
            response['status'] = 'Two or Three rolled, you won the Don\'t Pass bet'
        elif roll_sum == 12:
            response['status'] = 'Twelve rolled, it\'s a push'
        else:
            session['established_points'].append(roll_sum)
            if roll_sum in [4, 10]:
                session['current_bets'].append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 90 if session['initial_bankroll'] == 1000 else 100, 'potential_win': 45 if session['initial_bankroll'] == 1000 else 100})
            elif roll_sum in [5, 9]:
                session['current_bets'].append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 60 if session['initial_bankroll'] == 1000 else 75, 'potential_win': 40 if session['initial_bankroll'] == 1000 else 75})
            elif roll_sum in [6, 8]:
                session['current_bets'].append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 30 if session['initial_bankroll'] == 1000 else 60, 'potential_win': 25 if session['initial_bankroll'] == 1000 else 60})
            if len(session['established_points']) < 4:
                session['current_bets'].append({'type': 'Don\'t Come', 'bet': initial_bet, 'potential_win': initial_bet})
                session['bankroll'] -= initial_bet
            response['status'] = 'Point established'
    else:
        if roll_sum == 7:
            total_wins = sum(bet['potential_win'] for bet in session['current_bets'])
            session['bankroll'] += total_wins
            session['wins'] += total_wins
            response['status'] = 'Seven rolled, you won on all points'
            response['summary'] = {
                'initial_bankroll': session['initial_bankroll'],
                'final_bankroll': session['bankroll'],
                'roll_count': session['roll_count'],
                'results': session['results'],
                'points_established': len(session['established_points']),
                'points_hit': session['points_hit'],
                'wins': session['wins'],
                'losses': session['losses']
            }
        elif roll_sum in session['established_points']:
            session['bankroll'] -= initial_bet
            session['losses'] += initial_bet
            session['points_hit'] += 1
            response['status'] = 'Point hit, you lost a bet'
            if session['points_hit'] >= 3:
                response['status'] += '. Three strikes, wait for a 7 before betting again.'
                session['current_bets'] = []
        else:
            response['status'] = 'continue'

    response['bankroll'] = session['bankroll']
    response['wins'] = session['wins']
    response['losses'] = session['losses']

    return jsonify(response)

@app.route('/take_action', methods=['POST'])
def take_action():
    action = request.form['action']
    response = {
        'status': 'continue',
        'current_bets': session['current_bets'],
        'bankroll': session['bankroll'],
        'results': session['results'],
        'roll_count': session['roll_count'],
        'wins': session['wins'],
        'losses': session['losses']
    }

    if action == 'continue':
        roll_sum = game_roll()
        response['roll_sum'] = roll_sum
        initial_bet = 25 if session['initial_bankroll'] == 1000 else 50

        if roll_sum == 7:
            total_wins = sum(bet['potential_win'] for bet in session['current_bets'])
            session['bankroll'] += total_wins
            session['wins'] += total_wins
            response['status'] = 'Seven rolled, you won on all points'
            response['summary'] = {
                'initial_bankroll': session['initial_bankroll'],
                'final_bankroll': session['bankroll'],
                'roll_count': session['roll_count'],
                'results': session['results'],
                'points_established': len(session['established_points']),
                'points_hit': session['points_hit'],
                'wins': session['wins'],
                'losses': session['losses']
            }
        elif roll_sum in session['established_points']:
            session['bankroll'] -= initial_bet
            session['losses'] += initial_bet
            session['points_hit'] += 1
            response['status'] = 'Point hit, you lost a bet'
            if session['points_hit'] >= 3:
                response['status'] += '. Three strikes, wait for a 7 before betting again.'
                session['current_bets'] = []
        else:
            response['status'] = 'continue'
    elif action == 'no_bet':
        roll_sum = game_roll()
        response['roll_sum'] = roll_sum
        response['status'] = 'No bet'

    response['bankroll'] = session['bankroll']
    response['wins'] = session['wins']
    response['losses'] = session['losses']

    return jsonify(response)

@app.route('/summary', methods=['GET'])
def summary():
    return jsonify({
        'initial_bankroll': session.get('initial_bankroll'),
        'final_bankroll': session.get('bankroll'),
        'roll_count': session.get('roll_count'),
        'results': session.get('results'),
        'points_established': len(session.get('established_points', [])),
        'points_hit': session.get('points_hit', 0),
        'wins': session.get('wins', 0),
        'losses': session.get('losses', 0)
    })

if __name__ == '__main__':
    app.run(debug=True)
