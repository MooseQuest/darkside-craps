import random
import logging
import os
import sys
from dotenv import load_dotenv
from flask import Flask, render_template, request, jsonify, session
from flask_session import Session
from flask_socketio import SocketIO, emit
from pymongo import MongoClient
from bson.objectid import ObjectId
from redis import Redis

load_dotenv()  # Load environment variables from .env file

MONGO_URI = os.getenv('MONGO_URI')

client = MongoClient(MONGO_URI)
db = client.craps_game

# Set up logging
# logging.basicConfig(filename='craps_game.log', level=logging.INFO, 
#                     format='%(asctime)s - %(levelname)s - %(message)s')

# try:
# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
handler = logging.StreamHandler(sys.stdout)  # Use sys.stdout for Heroku logging
handler.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
handler.setFormatter(formatter)
logging.getLogger().addHandler(handler)

# Set up MongoDB logging
mongo_logger = logging.getLogger('pymongo')
mongo_logger.setLevel(logging.INFO)
mongo_logger.addHandler(handler)


# Flask Session and app config 
app = Flask(__name__)
app.config['SECRET_KEY'] = os.getenv('SECRET_KEY', 'default_secret_key')
app.config['SESSION_TYPE'] = 'redis'
app.config['SESSION_REDIS'] = Redis.from_url(os.getenv('REDISCLOUD_URL'))
Session(app)
socketio = SocketIO(app, cors_allowed_origins=["http://127.0.0.1:5000", "https://your-heroku-app.herokuapp.com"])

# Set up MongoDB connection
MONGO_URI = os.getenv('MONGO_URI')
client = MongoClient(MONGO_URI)
db = client.craps_game


# Function to simulate a dice roll
def roll_dice():
    return random.randint(1, 6), random.randint(1, 6)

class CrapsStateMachine:
    def __init__(self, initial_bankroll):
        self.initial_bankroll = initial_bankroll
        self.bankroll = initial_bankroll
        self.total_wagered = 0
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

        # Log the initialization
        logging.info(f'Initialized game with bankroll: {self.initial_bankroll}')

    # Sets up a new game and estalishes the initial state and the point off.
    def handle_event(self, event, bet_choice):
        state = 'come_out_roll' if self.come_out_roll else 'point_established'
        method_name = f'{state}_{event}'
        method = getattr(self, method_name, lambda: "Invalid state transition")
        return method(bet_choice)

    def come_out_roll_roll(self, bet_choice):
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        initial_bet = 25 if self.initial_bankroll == 1000 else 50
        lay_bet, lay_message = self._get_lay_bet(roll_sum)
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
            'strikes': self.strikes,
            'total_wagered': self.total_wagered,
            'next_bet': lay_message
        }

        if bet_choice == 'bet':
            self.total_wagered += initial_bet
            if roll_sum in [7, 11]:
                self.bankroll -= initial_bet
                self.losses += initial_bet
                response['status'] = 'Seven or Eleven rolled, you lost the Don\'t Pass bet'
                response['summary_button'] = True
                response['continue_button'] = True
                self.come_out_roll = True  # Reset to come-out roll
                logging.info(f'Loss on come out roll: {initial_bet}')
            elif roll_sum in [2, 3]:
                self.bankroll += initial_bet
                self.wins += initial_bet
                response['status'] = 'Two or Three rolled, you won the Don\'t Pass bet'
                logging.info(f'Win on come out roll: {initial_bet}')
            elif roll_sum == 12:
                response['status'] = 'Twelve rolled, it\'s a push'
                logging.info('Push on come out roll')
            else:
                if roll_sum not in [2, 3, 7, 11, 12]:
                    self.established_points.append(roll_sum)
                self.come_out_roll = False  # Not the come-out roll anymore
                self._place_odds_bet(roll_sum, initial_bet)
                response['status'] = 'Point established. Place a $' + str(initial_bet) + ' Don\'t Come bet.'
                logging.info(f'Point established: {roll_sum}')
        else:
            response['status'] = 'No bet'
        
        self.results.append(dice_roll)
        self.roll_count += 1
        logging.info(f'Roll: {dice_roll}, Roll Sum: {roll_sum}')
        return response

    def point_established_roll(self, bet_choice):
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        initial_bet = 25 if self.initial_bankroll == 1000 else 50
        lay_bet, lay_message = self._get_lay_bet(roll_sum)
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
            'strikes': self.strikes,
            'total_wagered': self.total_wagered,
            'next_bet': lay_message
        }

        if bet_choice == 'bet':
            self.total_wagered += initial_bet
            if roll_sum == 7:
                total_wins = sum(bet['potential_win'] for bet in self.current_bets if bet.get('point') in self.established_points)
                self.bankroll += total_wins
                self.wins += total_wins
                response['status'] = 'Seven rolled, you won on all points'
                response['summary_button'] = True
                response['continue_button'] = True
                self.come_out_roll = True  # Reset to come-out roll
                self.established_points = []  # Clear all established points
                logging.info(f'Win on established points: {total_wins}')
            elif roll_sum in self.established_points:
                point_bet = next((bet for bet in self.current_bets if bet.get('point') == roll_sum), None)
                self.bankroll -= (initial_bet + point_bet['bet'] if point_bet else 0)
                self.losses += (initial_bet + point_bet['bet'] if point_bet else 0)
                self.points_lost += 1
                self.established_points.remove(roll_sum)
                response['status'] = 'Point hit. Don\'t Behind ' + str(roll_sum) + ' and it\'s replaced by your Don\'t Come bet. This is strike ' + str(self.strikes + 1) + '. Place another Don\'t Come bet.'
                self.strikes += 1
                logging.info(f'Loss on established point: {roll_sum}')
                if self.strikes >= 3:
                    response['status'] += ' Three strikes, wait for a 7 before betting again.'
                    self.current_bets = []
            else:
                if roll_sum not in [2, 3, 7, 11, 12]:
                    self.established_points.append(roll_sum)
                self._place_odds_bet(roll_sum, initial_bet)
                response['status'] = 'Point established. Place a $' + str(initial_bet) + ' Don\'t Come bet.'
                logging.info(f'Point established: {roll_sum}')
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
                logging.info(f'Win on established points: {total_wins}')
            elif roll_sum in self.established_points:
                point_bet = next((bet for bet in self.current_bets if bet.get('point') == roll_sum), None)
                self.bankroll -= (initial_bet + point_bet['bet'] if point_bet else 0)
                self.losses += (initial_bet + point_bet['bet'] if point_bet else 0)
                self.points_lost += 1
                self.established_points.remove(roll_sum)
                response['status'] = 'Point hit. Don\'t Behind ' + str(roll_sum) + ' and it\'s replaced by your Don\'t Come bet. This is strike ' + str(self.strikes + 1) + '. Place another Don\'t Come bet.'
                self.strikes += 1
                logging.info(f'Loss on established point: {roll_sum}')
                if self.strikes >= 3:
                    response['status'] += ' Three strikes, wait for a 7 before betting again.'
                    self.current_bets = []
            else:
                if roll_sum not in [2, 3, 7, 11, 12]:
                    self.established_points.append(roll_sum)
                response['status'] = 'No bet'
        
        self.results.append(dice_roll)
        self.roll_count += 1
        logging.info(f'Roll: {dice_roll}, Roll Sum: {roll_sum}')
        return response

    def _place_odds_bet(self, roll_sum, initial_bet):
        if roll_sum in [4, 10]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 90 if self.initial_bankroll == 1000 else 100, 'potential_win': 45 if self.initial_bankroll == 1000 else 100})
        elif roll_sum in [5, 9]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 60 if self.initial_bankroll == 1000 else 75, 'potential_win': 40 if self.initial_bankroll == 1000 else 75})
        elif roll_sum in [6, 8]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 30 if self.initial_bankroll == 1000 else 60, 'potential_win': 25 if self.initial_bankroll == 1000 else 60})
        logging.info(f'Placed odds bet: Point {roll_sum}, Bet {initial_bet}')

    def _get_lay_bet(self, roll_sum):
        if roll_sum in [4, 10]:
            return (90 if self.initial_bankroll == 1000 else 100), 'Lay $90 if bankroll is $1000 or $100 if bankroll is $2000 on 4 or 10'
        elif roll_sum in [5, 9]:
            return (60 if self.initial_bankroll == 1000 else 75), 'Lay $60 if bankroll is $1000 or $75 if bankroll is $2000 on 5 or 9'
        elif roll_sum in [6, 8]:
            return (30 if self.initial_bankroll == 1000 else 60), 'Lay $30 if bankroll is $1000 or $60 if bankroll is $2000 on 6 or 8'
        else:
            return 0, 'No lay bet'

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
    socketio.emit('update_game_state', {
        'initial_bankroll': initial_bankroll,
        'current_bets': game.current_bets,
        'bankroll': game.bankroll,
        'total_wagered': game.total_wagered
    })
    logging.info(f'Game started with bankroll: {initial_bankroll}, session: {session}')
    return jsonify({
        'status': 'Game started',
        'initial_bankroll': initial_bankroll,
        'current_bets': game.current_bets,
        'bankroll': game.bankroll,
        'total_wagered': game.total_wagered
    })

@app.route('/roll', methods=['POST'])
def roll():
    if 'game_id' not in session:
        return jsonify({'error': 'No active game'}), 400
    
    game_id = ObjectId(session['game_id'])
    game_data = db.games.find_one({'_id': game_id})
    game = CrapsStateMachine(game_data['initial_bankroll'])
    game.__dict__.update(game_data)
    bet_choice = request.form['bet_choice']
    response = game.handle_event('roll', bet_choice)
    db.games.update_one({'_id': game_id}, {'$set': game.__dict__})
    socketio.emit('update_game_state', response)
    logging.info(f'Roll executed with bet choice: {bet_choice}, session: {session}')
    return jsonify(response)

@app.route('/summary', methods=['GET'])
def summary():
    if 'game_id' not in session:
        return jsonify({'error': 'No active game'}), 400
    
    game_id = ObjectId(session['game_id'])
    game_data = db.games.find_one({'_id': game_id})
    
    summary = {
        'initial_bankroll': game_data['initial_bankroll'],
        'final_bankroll': game_data['bankroll'],
        'roll_count': game_data['roll_count'],
        'results': game_data['results'],
        'points_established': len(game_data['established_points']),
        'points_hit': game_data['points_hit'],
        'points_lost': game_data['points_lost'],
        'wins': game_data['wins'],
        'losses': game_data['losses'],
        'strikes': game_data['strikes'],
        'total_wagered': game_data['total_wagered']
    }
    socketio.emit('update_game_summary', summary)
    logging.info(f'Game summary: {summary}')
    return jsonify(summary)

if __name__ == '__main__':
    if os.environ.get('FLASK_ENV') == 'development':
        port = int(os.environ.get("PORT", 5000))
        socketio.run(app, host='0.0.0.0', port=port, debug=True)
    else:
        port = int(os.environ.get("PORT", 5000))
        socketio.run(app, host='0.0.0.0', port=port, allow_unsafe_werkzeug=True)