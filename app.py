import random
from flask import Flask, render_template, request, jsonify, session
from flask_session import Session
from pymongo import MongoClient
from bson.objectid import ObjectId

app = Flask(__name__)
app.config['SECRET_KEY'] = 'supersecretkey'
app.config['SESSION_TYPE'] = 'filesystem'
Session(app)

client = MongoClient('mongodb://localhost:27017/')
db = client.craps_game

# Function to simulate a dice roll
def roll_dice():
    return random.randint(1, 6), random.randint(1, 6)

# Initialize game state
def init_game(initial_bankroll):
    game_data = {
        'initial_bankroll': initial_bankroll,
        'bankroll': initial_bankroll,
        'current_bets': [{'type': 'Don\'t Pass', 'bet': 25 if initial_bankroll == 1000 else 50, 'potential_win': 25 if initial_bankroll == 1000 else 50}],
        'established_points': [],
        'results': [],
        'roll_count': 0,
        'points_hit': 0,
        'wins': 0,
        'losses': 0,
        'strikes': 0
    }
    game_id = db.games.insert_one(game_data).inserted_id
    return game_id

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/start', methods=['POST'])
def start_game():
    initial_bankroll = int(request.form['bankroll'])
    game_id = init_game(initial_bankroll)
    game_data = db.games.find_one({'_id': game_id})
    session['game_id'] = str(game_id)
    session['initial_bankroll'] = initial_bankroll
    return jsonify({
        'status': 'Game started',
        'initial_bankroll': initial_bankroll,
        'current_bets': game_data['current_bets'],
        'bankroll': game_data['bankroll']
    })

@app.route('/roll', methods=['POST'])
def roll():
    game_id = ObjectId(session['game_id'])
    game_data = db.games.find_one({'_id': game_id})
    
    dice_roll = roll_dice()
    roll_sum = sum(dice_roll)
    initial_bet = 25 if game_data['initial_bankroll'] == 1000 else 50
    response = {
        'status': 'continue',
        'roll_sum': roll_sum,
        'dice_roll': dice_roll,
        'current_bets': game_data['current_bets'],
        'bankroll': game_data['bankroll'],
        'results': game_data['results'],
        'roll_count': game_data['roll_count'],
        'wins': game_data['wins'],
        'losses': game_data['losses']
    }

    if len(game_data['established_points']) < 4:
        if roll_sum in [7, 11]:
            game_data['bankroll'] -= initial_bet
            game_data['losses'] += initial_bet
            response['status'] = 'Seven or Eleven rolled, you lost the Don\'t Pass bet'
        elif roll_sum in [2, 3]:
            game_data['bankroll'] += initial_bet
            game_data['wins'] += initial_bet
            response['status'] = 'Two or Three rolled, you won the Don\'t Pass bet'
        elif roll_sum == 12:
            response['status'] = 'Twelve rolled, it\'s a push'
        else:
            game_data['established_points'].append(roll_sum)
            if roll_sum in [4, 10]:
                game_data['current_bets'].append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 90 if game_data['initial_bankroll'] == 1000 else 100, 'potential_win': 45 if game_data['initial_bankroll'] == 1000 else 100})
            elif roll_sum in [5, 9]:
                game_data['current_bets'].append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 60 if game_data['initial_bankroll'] == 1000 else 75, 'potential_win': 40 if game_data['initial_bankroll'] == 1000 else 75})
            elif roll_sum in [6, 8]:
                game_data['current_bets'].append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 30 if game_data['initial_bankroll'] == 1000 else 60, 'potential_win': 25 if game_data['initial_bankroll'] == 1000 else 60})
            if len(game_data['established_points']) < 4:
                game_data['current_bets'].append({'type': 'Don\'t Come', 'bet': initial_bet, 'potential_win': initial_bet})
                game_data['bankroll'] -= initial_bet
            response['status'] = 'Point established'
    else:
        if roll_sum == 7:
            total_wins = sum(bet['potential_win'] for bet in game_data['current_bets'])
            game_data['bankroll'] += total_wins
            game_data['wins'] += total_wins
            response['status'] = 'Seven rolled, you won on all points'
            response['summary'] = {
                'initial_bankroll': game_data['initial_bankroll'],
                'final_bankroll': game_data['bankroll'],
                'roll_count': game_data['roll_count'],
                'results': game_data['results'],
                'points_established': len(game_data['established_points']),
                'points_hit': game_data['points_hit'],
                'wins': game_data['wins'],
                'losses': game_data['losses']
            }
        elif roll_sum in game_data['established_points']:
            game_data['bankroll'] -= initial_bet
            game_data['losses'] += initial_bet
            game_data['points_hit'] += 1
            response['status'] = 'Point hit, you lost a bet'
            if game_data['points_hit'] >= 3:
                response['status'] += '. Three strikes, wait for a 7 before betting again.'
                game_data['current_bets'] = []
        else:
            response['status'] = 'continue'

    game_data['results'].append(dice_roll)
    game_data['roll_count'] += 1
    db.games.update_one({'_id': game_id}, {'$set': game_data})

    response['bankroll'] = game_data['bankroll']
    response['wins'] = game_data['wins']
    response['losses'] = game_data['losses']

    return jsonify(response)

@app.route('/take_action', methods=['POST'])
def take_action():
    action = request.form['action']
    game_id = ObjectId(session['game_id'])
    game_data = db.games.find_one({'_id': game_id})
    
    response = {
        'status': 'continue',
        'current_bets': game_data['current_bets'],
        'bankroll': game_data['bankroll'],
        'results': game_data['results'],
        'roll_count': game_data['roll_count'],
        'wins': game_data['wins'],
        'losses': game_data['losses']
    }

    if action == 'continue':
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        response['roll_sum'] = roll_sum
        response['dice_roll'] = dice_roll
        initial_bet = 25 if game_data['initial_bankroll'] == 1000 else 50

        if roll_sum == 7:
            total_wins = sum(bet['potential_win'] for bet in game_data['current_bets'])
            game_data['bankroll'] += total_wins
            game_data['wins'] += total_wins
            response['status'] = 'Seven rolled, you won on all points'
            response['summary'] = {
                'initial_bankroll': game_data['initial_bankroll'],
                'final_bankroll': game_data['bankroll'],
                'roll_count': game_data['roll_count'],
                'results': game_data['results'],
                'points_established': len(game_data['established_points']),
                'points_hit': game_data['points_hit'],
                'wins': game_data['wins'],
                'losses': game_data['losses']
            }
        elif roll_sum in game_data['established_points']:
            game_data['bankroll'] -= initial_bet
            game_data['losses'] += initial_bet
            game_data['points_hit'] += 1
            response['status'] = 'Point hit, you lost a bet'
            if game_data['points_hit'] >= 3:
                response['status'] += '. Three strikes, wait for a 7 before betting again.'
                game_data['current_bets'] = []
        else:
            response['status'] = 'continue'

        game_data['results'].append(dice_roll)
        game_data['roll_count'] += 1
        db.games.update_one({'_id': game_id}, {'$set': game_data})

    elif action == 'no_bet':
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        response['roll_sum'] = roll_sum
        response['dice_roll'] = dice_roll
        response['status'] = 'No bet'
        game_data['results'].append(dice_roll)
        game_data['roll_count'] += 1
        db.games.update_one({'_id': game_id}, {'$set': game_data})

    response['bankroll'] = game_data['bankroll']
    response['wins'] = game_data['wins']
    response['losses'] = game_data['losses']

    return jsonify(response)

@app.route('/summary', methods=['GET'])
def summary():
    game_id = ObjectId(session['game_id'])
    game_data = db.games.find_one({'_id': game_id})
    
    return jsonify({
        'initial_bankroll': game_data['initial_bankroll'],
        'final_bankroll': game_data['bankroll'],
        'roll_count': game_data['roll_count'],
        'results': game_data['results'],
        'points_established': len(game_data['established_points']),
        'points_hit': game_data['points_hit'],
        'wins': game_data['wins'],
        'losses': game_data['losses']
    })

if __name__ == '__main__':
    app.run(debug=True)
