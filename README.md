
# Dark Side Craps Strategy 2024

This repository contains the initial setup for developing a craps strategy model using Python and Flask. The objective of this model is to simulate a craps strategy and determine the likelihood of success.

This hopes to illustrate and run the Long Dark Side Strategy that I have put together to ensure you have a long play at the crap table. The idea is not to win or lose, but play close to the house edge and for a long period of time. The likelihood of winning at this strategy is low, but higher than if you played the light side. Thus, there are larger bankrolls involved.

Odds here are maintained around the idea of betting enough to equal or gain higher than your original bet. You may get questions about your bet, but nonetheless this has been tested.

Discipline is required to execute this strategy as it is built to ensure you don't lose all the money faster than a certain amount of play time. As the time extended for your play, along with the number of roles, the probably of success goes down. The limitations of that are built into the simulation to let you know when it's time to walk away.

## Folder Structure

```
.
├── .gitignore
├── .vscode
│   └── launch.json
├── README.md
├── app.py
├── docs
│   ├── KH-Long-Dark-Side.md
│   └── chatgpt-prompt.txt
├── templates
│   └── index.html
├── test_craps_game.py
├── test_flask_app.py
└── test_strategy.py
```

## Prerequisites

To run this project, you will need to have the following software installed on your local machine:

- Python 3.x
- pip (Python package installer)

## ChatGPT Prompt

If you want to journey down creating or building this app with generative AI, here is the prompt I used along with the Strategy. First I wrote out the strategy, the built out the components to support the strategy.

## Setup Instructions

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd something-something-darkside
   ```

2. **Create and Activate a Virtual Environment**

   On macOS/Linux:
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   ```

   On Windows:
   ```bash
   python -m venv venv
   .\venv\Scripts\activate
   ```

3. **Install Required Packages**

   ```bash
   pip install Flask
   ```

## Running the Application

1. **Start the Flask Application**

   ```bash
   python app.py
   ```

2. **Access the Application**

   Open your web browser and navigate to `http://127.0.0.1:5000` to view and interact with the craps strategy simulation.

## Application Structure

- `app.py`: The main Flask application file. It contains the routes and logic for the craps simulation.
- `docs/`: This directory contains documentation files.
  - `KH-Long-Dark-Side.md`: Detailed description of the KH Long Play Don’t Strategy.
  - `chatgpt-prompt.txt`: The initial prompt used to generate the project files.
- `templates/`: This directory contains the HTML template for the web interface.
  - `index.html`: The main HTML file for the Flask application.

## How to Use

1. **Open the Application**:
   Navigate to `http://127.0.0.1:5000` in your web browser.

2. **Enter Initial Bankroll**:
   Enter the initial bankroll amount in the provided input field.

3. **Run Simulation**:
   Click on the "Run Simulation" button to start the craps strategy simulation.

4. **View Results**:
   The simulation results will be displayed below the form, showing the changes in bankroll over time based on the craps strategy.

## Contributing

If you would like to contribute to this project, please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more information.
