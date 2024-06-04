
# Craps Strategy 2024

This repository contains the initial setup for developing a craps strategy model using Python and Flask. The objective of this model is to simulate a craps strategy and determine the likelihood of success.

## Folder Structure

```
something-something-darkside/
├── README.md
├── app.py
├── docs/
│   ├── KH-Long-Dark-Side.md
│   └── chatgpt-prompt.txt
└── templates/
    └── index.html
```

## Prerequisites

To run this project, you will need to have the following software installed on your local machine:

- Python 3.x
- pip (Python package installer)

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
