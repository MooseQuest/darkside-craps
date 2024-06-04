import random
from flask import Flask, render_template, request, jsonify

app = Flask(__name__)

# Function to simulate a dice roll
def roll_dice():
    return random.randint(1, 6), random.randint(1, 6)

# Function to simulate a single game of craps
def simulate_game(initial_bankroll):
    bankroll = initial_bankroll
    results = []
    
    while bankroll > 0:
        # Come out roll
        roll = roll_dice()
        point = sum(roll)
        
        if point in [7, 11]:
            bankroll -= 25
        elif point in [2, 3, 12]:
            bankroll += 25
        else:
            # Point established
            while True:
                roll = roll_dice()
                roll_sum = sum(roll)
                if roll_sum == point:
                    bankroll -= 25
                    break
                elif roll_sum == 7:
                    bankroll += 25
                    break
        
        results.append(bankroll)
    
    return results

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/simulate', methods=['POST'])
def simulate():
    try:
        initial_bankroll = int(request.form['bankroll'])
        results = simulate_game(initial_bankroll)
        return jsonify(results)
    except Exception as e:
        return str(e), 500

if __name__ == '__main__':
    app.run(debug=True)
