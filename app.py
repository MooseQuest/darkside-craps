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

class CrapsStateMachine:
    def __init__(self, initial_bankroll):
        self.initial_bankroll = initial_bankroll
        self.bankroll = initial_bankroll
        self.current_bets = [{'type': 'Don\'t Pass', 'bet': 25 if initial_bankroll == 1000 else 50, 'potential_win': 25 if initial_bankroll == 1000 else 50}]
        self.established_points = []
        self.results = []
        self.roll_count = 0
        self.points_hit = 0
        self.points_lost = 0
        self.wins = 0
        self.losses = 0
        self.strikes = 0
        self.come_out_roll = True

    def handle_event(self, event, bet_choice):
        state = 'come_out_roll' if self.come_out_roll else 'point_established'
        method_name = f'{state}_{event}'
        method = getattr(self, method_name, lambda: "Invalid state transition")
        return method(bet_choice)

    def come_out_roll_roll(self, bet_choice):
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        initial_bet = 25 if self.initial_bankroll == 1000 else 50
        response = {
            'status': 'continue',
            'roll_sum': roll_sum,
            'dice_roll': dice_roll,
            'current_bets': self.current_bets,
            'bankroll': self.bankroll,
            'results': self.results,
            'roll_count': self.roll_count,
            'wins': self.wins,
            'losses': self.losses,
            'established_points': self.established_points,
            'strikes': self.strikes
        }

        if bet_choice == 'bet':
            if roll_sum in [7, 11]:
                self.bankroll -= initial_bet
                self.losses += initial_bet
                response['status'] = 'Seven or Eleven rolled, you lost the Don\'t Pass bet'
                response['summary_button'] = True
                response['continue_button'] = True
                self.come_out_roll = True  # Reset to come-out roll
            elif roll_sum in [2, 3]:
                self.bankroll += initial_bet
                self.wins += initial_bet
                response['status'] = 'Two or Three rolled, you won the Don\'t Pass bet'
            elif roll_sum == 12:
                response['status'] = 'Twelve rolled, it\'s a push'
            else:
                self.established_points.append(roll_sum)
                self.come_out_roll = False  # Not the come-out roll anymore
                self._place_odds_bet(roll_sum, initial_bet)
                response['status'] = 'Point established. Place a $' + str(initial_bet) + ' Don\'t Come bet.'
        else:
            response['status'] = 'No bet'
        
        self.results.append(dice_roll)
        self.roll_count += 1
        return response

    def point_established_roll(self, bet_choice):
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        initial_bet = 25 if self.initial_bankroll == 1000 else 50
        response = {
            'status': 'continue',
            'roll_sum': roll_sum,
            'dice_roll': dice_roll,
            'current_bets': self.current_bets,
            'bankroll': self.bankroll,
            'results': self.results,
            'roll_count': self.roll_count,
            'wins': self.wins,
            'losses': self.losses,
            'established_points': self.established_points,
            'strikes': self.strikes
        }

        if bet_choice == 'bet':
            if roll_sum == 7:
                total_wins = sum(bet['potential_win'] for bet in self.current_bets if bet.get('point') in self.established_points)
                self.bankroll += total_wins
                self.wins += total_wins
                response['status'] = 'Seven rolled, you won on all points'
                response['summary_button'] = True
                response['continue_button'] = True
                self.come_out_roll = True  # Reset to come-out roll
                self.established_points = []  # Clear all established points
            elif roll_sum in self.established_points:
                point_bet = next((bet for bet in self.current_bets if bet.get('point') == roll_sum), None)
                self.bankroll -= (initial_bet + point_bet['bet'] if point_bet else 0)
                self.losses += (initial_bet + point_bet['bet'] if point_bet else 0)
                self.points_lost += 1
                self.established_points.remove(roll_sum)
                response['status'] = 'Point hit. Don\'t Behind ' + str(roll_sum) + ' and it\'s replaced by your Don\'t Come bet. This is strike ' + str(self.strikes + 1) + '. Place another Don\'t Come bet.'
                self.strikes += 1
                if self.strikes >= 3:
                    response['status'] += ' Three strikes, wait for a 7 before betting again.'
                    self.current_bets = []
            else:
                self._place_odds_bet(roll_sum, initial_bet)
                response['status'] = 'Point established. Place a $' + str(initial_bet) + ' Don\'t Come bet.'
        else:
            if roll_sum == 7:
                total_wins = sum(bet['potential_win'] for bet in self.current_bets if bet.get('point') in self.established_points)
                self.bankroll += total_wins
                self.wins += total_wins
                response['status'] = 'Seven rolled, you won on all points'
                response['summary_button'] = True
                response['continue_button'] = True
                self.come_out_roll = True  # Reset to come-out roll
                self.established_points = []  # Clear all established points
            elif roll_sum in self.established_points:
                point_bet = next((bet for bet in self.current_bets if bet.get('point') == roll_sum), None)
                self.bankroll -= (initial_bet + point_bet['bet'] if point_bet else 0)
                self.losses += (initial_bet + point_bet['bet'] if point_bet else 0)
                self.points_lost += 1
                self.established_points.remove(roll_sum)
                response['status'] = 'Point hit. Don\'t Behind ' + str(roll_sum) + ' and it\'s replaced by your Don\'t Come bet. This is strike ' + str(self.strikes + 1) + '. Place another Don\'t Come bet.'
                self.strikes += 1
                if self.strikes >= 3:
                    response['status'] += ' Three strikes, wait for a 7 before betting again.'
                    self.current_bets = []
            else:
                response['status'] = 'No bet'
        
        self.results.append(dice_roll)
        self.roll_count += 1
        return response

    def _place_odds_bet(self, roll_sum, initial_bet):
        if roll_sum in [4, 10]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 90 if self.initial_bankroll == 1000 else 100, 'potential_win': 45 if self.initial_bankroll == 1000 else 100})
        elif roll_sum in [5, 9]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 60 if self.initial_bankroll == 1000 else 75, 'potential_win': 40 if self.initial_bankroll == 1000 else 75})
        elif roll_sum in [6, 8]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 30 if self.initial_bankroll == 1000 else 60, 'potential_win': 25 if self.initial_bankroll == 1000 else 60})

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/start', methods=['POST'])
def start_game():
    initial_bankroll = int(request.form['bankroll'])
    game = CrapsStateMachine(initial_bankroll)
    game_id = db.games.insert_one(game.__dict__).inserted_id
    session['game_id'] = str(game_id)
    session['initial_bankroll'] = initial_bankroll
    return jsonify({
        'status': 'Game started',
        'initial_bankroll': initial_bankroll,
        'current_bets': game.current_bets,
        'bankroll': game.bankroll
    })

@app.route('/roll', methods=['POST'])
def roll():
    game_id = ObjectId(session['game_id'])
    game_data = db.games.find_one({'_id': game_id})
    game = CrapsStateMachine(game_data['initial_bankroll'])
    game.__dict__.update(game_data)
    bet_choice = request.form['bet_choice']
    response = game.handle_event('roll', bet_choice)
    db.games.update_one({'_id': game_id}, {'$set': game.__dict__})
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
        'points_lost': game_data['points_lost'],
        'wins': game_data['wins'],
        'losses': game_data['losses'],
        'strikes': game_data['strikes']
    })

if __name__ == '__main__':
    app.run(debug=True)
