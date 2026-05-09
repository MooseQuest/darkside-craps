import json
import os
import sys
import time

import requests

BASE_URL = os.environ.get('CRAPS_URL', 'http://127.0.0.1:5000')


def start_game(session, bankroll, log_file='bet_tracker.log'):
    response = session.post(f'{BASE_URL}/start', data={'bankroll': bankroll})
    response.raise_for_status()
    return response.json()


def roll_dice(session, bet_choice, log_file):
    response = session.post(f'{BASE_URL}/roll', data={'bet_choice': bet_choice})
    response.raise_for_status()
    return response.json()


def get_summary(session, log_file):
    response = session.get(f'{BASE_URL}/summary')
    response.raise_for_status()
    return response.json()

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
    failures = []

    session = requests.Session()

    try:
        game_data = start_game(session, initial_bankroll)
        print("Game started:", game_data)
        log_result(log_file, f"Game started: {json.dumps(game_data)}")
        if game_data.get('initial_bankroll') != initial_bankroll:
            failures.append(f"start: initial_bankroll mismatch ({game_data.get('initial_bankroll')!r})")
    except Exception as e:
        print(f"FAIL start_game: {e}", file=sys.stderr)
        sys.exit(1)

    for i in range(10):
        try:
            response = roll_dice(session, 'bet', log_file)
            print(f"Roll {i+1}:", response)
            log_result(log_file, f"Roll {i+1}: {json.dumps(response)}")
            check_for_anomalies(response, log_file)
            if 'roll_sum' not in response or 'bankroll' not in response:
                failures.append(f"roll {i+1}: response missing required keys")
        except Exception as e:
            print(f"FAIL roll {i+1}: {e}", file=sys.stderr)
            sys.exit(1)
        time.sleep(0.2)

    try:
        summary = get_summary(session, log_file)
        print("Game summary:", summary)
        log_result(log_file, f"Game summary: {json.dumps(summary)}")
        for required_key in ('initial_bankroll', 'final_bankroll', 'roll_count', 'total_wagered'):
            if required_key not in summary:
                failures.append(f"summary: missing {required_key!r}")
        if summary.get('roll_count', 0) < 10:
            failures.append(f"summary: roll_count={summary.get('roll_count')} (expected >=10)")
    except Exception as e:
        print(f"FAIL get_summary: {e}", file=sys.stderr)
        sys.exit(1)

    if failures:
        print("\nTEST FAILURES:", file=sys.stderr)
        for f in failures:
            print(f"  - {f}", file=sys.stderr)
        sys.exit(1)
    print("\nAll smoke checks passed.")


if __name__ == '__main__':
    main()
