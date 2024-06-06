import requests
import json
import time

BASE_URL = 'http://127.0.0.1:5000'

def start_game(bankroll, log_file):  # Add log_file parameter
    try:
        response = requests.post(f'{BASE_URL}/start', data={'bankroll': bankroll})
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        print(f"Error starting game: {e}")
        log_result(log_file, f"Error starting game: {e}")
        return None

def roll_dice(bet_choice, log_file):  # Add log_file parameter
    try:
        response = requests.post(f'{BASE_URL}/roll', data={'bet_choice': bet_choice})
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        print(f"Error rolling dice: {e}")
        log_result(log_file, f"Error rolling dice: {e}")
        return None

def get_summary(log_file):  # Add log_file parameter
    try:
        response = requests.get(f'{BASE_URL}/summary')
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        print(f"Error getting summary: {e}")
        log_result(log_file, f"Error getting summary: {e}")
        return None

def log_result(log_file, message):
    with open(log_file, 'a') as file:
        file.write(message + '\n')

def check_for_anomalies(response, log_file):
    established_points = response.get('established_points', [])
    if 7 in established_points or 11 in established_points:
        anomaly_message = f"Anomaly detected: 7 or 11 in established points: {established_points}"
        print(anomaly_message)
        log_result(log_file, anomaly_message)

def main():
    log_file = 'bet_tracker.log'
    initial_bankroll = 1000

    # Start the game
    game_data = start_game(initial_bankroll, log_file)
    if game_data:
        print("Game started:", game_data)
        log_result(log_file, f"Game started: {json.dumps(game_data)}")
    else:
        return

    # Sequence of dice rolls
    dice_sequence = [
        ('bet', (6, 2)),  # Roll 8
        ('bet', (4, 2)),  # Roll 6
        ('bet', (1, 1)),  # Roll 2
        ('bet', (3, 4)),  # Roll 7
        ('bet', (5, 1)),  # Roll 6
        ('bet', (6, 3)),  # Roll 9
        ('bet', (5, 4)),  # Roll 9
        ('bet', (3, 3)),  # Roll 6
        ('bet', (2, 5)),  # Roll 7
        ('bet', (4, 4))   # Roll 8
    ]

    for bet_choice, dice in dice_sequence:
        response = roll_dice(bet_choice, log_file)
        if response:
            print("Roll response:", response)
            log_result(log_file, f"Roll response: {json.dumps(response)}")
            check_for_anomalies(response, log_file)
        else:
            break

        # Delay to simulate time between rolls
        time.sleep(1)

    # Get summary
    summary = get_summary(log_file)
    if summary:
        print("Game summary:", summary)
        log_result(log_file, f"Game summary: {json.dumps(summary)}")
    else:
        print("Failed to get game summary.")

if __name__ == '__main__':
    main()
