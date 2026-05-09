import gevent.monkey
gevent.monkey.patch_all()

import logging
import os
import random
import sys
import uuid

from bson.objectid import ObjectId
from dotenv import load_dotenv
from flask import Flask, jsonify, make_response, render_template, request, session
from flask_cors import CORS
from flask_session import Session
from flask_socketio import SocketIO
from pymongo import MongoClient

import secrets_loader

load_dotenv()

MONGO_URI = secrets_loader.load('MONGO_URI')
SESSION_MONGO_URI = secrets_loader.load('SESSION_MONGO_URI', default=MONGO_URI)

if not MONGO_URI:
    raise RuntimeError(
        "MONGO_URI is not set. Configure Infisical (INFISICAL_CLIENT_ID + "
        "INFISICAL_CLIENT_SECRET + INFISICAL_PROJECT_ID) or set MONGO_URI in env."
    )

client = MongoClient(MONGO_URI)
session_client = MongoClient(SESSION_MONGO_URI)
db = client.craps_game

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
CORS(app)
app.config['SECRET_KEY'] = secrets_loader.load('SECRET_KEY', default='default_secret_key')
app.config['SESSION_TYPE'] = 'mongodb'
app.config['SESSION_MONGODB'] = session_client
app.config['SESSION_MONGODB_DB'] = 'craps_game_sessions'
app.config['SESSION_MONGODB_COLLECT'] = 'sessions'
Session(app)
socketio = SocketIO(app, cors_allowed_origins="*", async_mode='gevent')

def generate_session_id():
    return str(uuid.uuid4())


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
        lay_bet, lay_message = self._get_lay_bet(roll_sum, on_come_out=True)
        response = {
            'status': 'continue',
            'roll_sum': roll_sum,
            'dice_roll': dice_roll,
            'next_bet': lay_message,
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
                self.current_bets = [{'type': 'Don\'t Pass', 'bet': initial_bet, 'potential_win': initial_bet}]
                logging.info(f'Loss on come out roll: {initial_bet}')
            elif roll_sum in [2, 3]:
                self.bankroll += initial_bet
                self.wins += initial_bet
                response['status'] = 'Two or Three rolled, you won the Don\'t Pass bet'
                response['summary_button'] = True
                self.current_bets = [{'type': 'Don\'t Pass', 'bet': initial_bet, 'potential_win': initial_bet}]
                logging.info(f'Win on come out roll: {initial_bet}')
            elif roll_sum == 12:
                response['status'] = 'Twelve rolled, it\'s a push'
                logging.info('Push on come out roll')
            else:
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
        self._snapshot_into(response)
        return response

    def point_established_roll(self, bet_choice):
        dice_roll = roll_dice()
        roll_sum = sum(dice_roll)
        initial_bet = 25 if self.initial_bankroll == 1000 else 50
        lay_bet, lay_message = self._get_lay_bet(roll_sum, on_come_out=False)
        response = {
            'status': 'continue',
            'roll_sum': roll_sum,
            'dice_roll': dice_roll,
            'next_bet': lay_message,
        }
        placing_bet = (bet_choice == 'bet')

        if placing_bet:
            self.total_wagered += initial_bet

        POINT_NUMBERS = (4, 5, 6, 8, 9, 10)

        if roll_sum == 7:
            self._resolve_seven_out(initial_bet, response)
        elif roll_sum in self.established_points:
            # Don't Come bet (initial_bet) only loses if a new one was placed this roll
            extra_loss = initial_bet if placing_bet else 0
            self._resolve_point_hit(roll_sum, extra_loss, response)
        elif roll_sum in POINT_NUMBERS and placing_bet:
            self.established_points.append(roll_sum)
            self._place_odds_bet(roll_sum, initial_bet)
            response['status'] = 'Point established. Place a $' + str(initial_bet) + ' Don\'t Come bet.'
            logging.info(f'Point established: {roll_sum}')
        else:
            # 2, 3, 11, 12 (or no-bet on a fresh point) — no effect on existing bets
            response['status'] = 'No effect on existing bets' if placing_bet else 'No bet'

        self.results.append(dice_roll)
        self.roll_count += 1
        logging.info(f'Roll: {dice_roll}, Roll Sum: {roll_sum}')
        self._snapshot_into(response)
        return response

    def _resolve_seven_out(self, initial_bet, response):
        total_wins = sum(bet['potential_win'] for bet in self.current_bets if bet.get('point') in self.established_points)
        self.bankroll += total_wins
        self.wins += total_wins
        response['status'] = 'Seven rolled, you won on all points'
        response['summary_button'] = True
        response['continue_button'] = True
        self.come_out_roll = True
        self.established_points = []
        self.current_bets = [{'type': 'Don\'t Pass', 'bet': initial_bet, 'potential_win': initial_bet}]
        self.strikes = 0
        logging.info(f'Win on established points: {total_wins}')

    def _resolve_point_hit(self, roll_sum, extra_loss, response):
        point_bet = next((bet for bet in self.current_bets if bet.get('point') == roll_sum), None)
        loss = extra_loss + (point_bet['bet'] if point_bet else 0)
        self.bankroll -= loss
        self.losses += loss
        self.points_hit += 1
        self.established_points.remove(roll_sum)
        if point_bet is not None:
            self.current_bets.remove(point_bet)
        response['status'] = 'Point hit. Don\'t Behind ' + str(roll_sum) + ' and it\'s replaced by your Don\'t Come bet. This is strike ' + str(self.strikes + 1) + '. Place another Don\'t Come bet.'
        self.strikes += 1
        logging.info(f'Loss on established point: {roll_sum}')
        if self.strikes >= 3:
            response['status'] += ' Three strikes, wait for a 7 before betting again.'
            self.current_bets = []

    def _snapshot_into(self, response):
        response['current_bets'] = list(self.current_bets)
        response['bankroll'] = self.bankroll
        response['results'] = list(self.results)
        response['roll_count'] = self.roll_count
        response['wins'] = self.wins
        response['losses'] = self.losses
        response['established_points'] = list(self.established_points)
        response['strikes'] = self.strikes
        response['total_wagered'] = self.total_wagered

    def _place_odds_bet(self, roll_sum, initial_bet):
        if roll_sum in [4, 10]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 90 if self.initial_bankroll == 1000 else 100, 'potential_win': 45 if self.initial_bankroll == 1000 else 100})
        elif roll_sum in [5, 9]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 60 if self.initial_bankroll == 1000 else 75, 'potential_win': 40 if self.initial_bankroll == 1000 else 75})
        elif roll_sum in [6, 8]:
            self.current_bets.append({'type': 'Lay Odds', 'point': roll_sum, 'bet': 30 if self.initial_bankroll == 1000 else 60, 'potential_win': 25 if self.initial_bankroll == 1000 else 60})
        logging.info(f'Placed odds bet: Point {roll_sum}, Bet {initial_bet}')

    def _get_lay_bet(self, roll_sum, on_come_out=False):
        # Lay-odds bets only apply when a new point is being established.
        # Suppress messaging when re-rolling on a 7/11/2/3/12 come-out or
        # when a 7 / repeated point resolves the round.
        if roll_sum not in [4, 5, 6, 8, 9, 10]:
            return 0, ''
        # If we're already past come-out and this number is already an
        # established point, no new lay-odds is placed.
        if not on_come_out and roll_sum in self.established_points:
            return 0, ''
        if roll_sum in [4, 10]:
            return (90 if self.initial_bankroll == 1000 else 100), 'Lay $90 if bankroll is $1000 or $100 if bankroll is $2000 on 4 or 10'
        if roll_sum in [5, 9]:
            return (60 if self.initial_bankroll == 1000 else 75), 'Lay $60 if bankroll is $1000 or $75 if bankroll is $2000 on 5 or 9'
        return (30 if self.initial_bankroll == 1000 else 60), 'Lay $30 if bankroll is $1000 or $60 if bankroll is $2000 on 6 or 8'

@app.route('/')
def home():
    session_id = request.cookies.get('session_id')
    if not session_id:
        session_id = generate_session_id()
        resp = make_response(render_template('index.html'))
        resp.set_cookie('session_id', session_id)
        return resp
    return render_template('index.html')

@app.route('/start', methods=['POST'])
def start_game():
    initial_bankroll = int(request.form['bankroll'])
    game = CrapsStateMachine(initial_bankroll)
    game_id = db.games.insert_one(game.__dict__).inserted_id
    session['game_id'] = str(game_id)
    session['initial_bankroll'] = initial_bankroll
    logging.info(f'Game started with bankroll: {initial_bankroll}, session: {session}, game_id: {game_id}')
    return jsonify({
        'status': 'Game started',
        'initial_bankroll': initial_bankroll,
        'current_bets': game.current_bets,
        'bankroll': game.bankroll,
        'total_wagered': game.total_wagered,
        'session_id': session.sid
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
    logging.info(f'Roll executed with bet choice: {bet_choice}, session: {session}, game_id: {game_id}')
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
    logging.info(f'Game summary: {summary}, session: {session}, game_id: {game_id}')
    return jsonify(summary)

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 5000))
    if os.environ.get('FLASK_ENV') == 'development':
        socketio.run(app, host='0.0.0.0', port=port, debug=True)
    else:
        from gevent import pywsgi
        from geventwebsocket.handler import WebSocketHandler
        server = pywsgi.WSGIServer(('', port), app, handler_class=WebSocketHandler)
        server.serve_forever()