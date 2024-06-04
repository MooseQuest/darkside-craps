from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
import time

# Path to your webdriver executable
# For example, if you're using Chrome, the path should point to your chromedriver
webdriver_path = '/path/to/chromedriver'

# Initialize the WebDriver
driver = webdriver.Chrome(executable_path=webdriver_path)

try:
    # Open the Flask application
    driver.get("http://127.0.0.1:5000")

    # Find the input element for the bankroll
    bankroll_input = driver.find_element(By.ID, "bankroll")
    
    # Clear the input field just in case there's any pre-filled text
    bankroll_input.clear()
    
    # Enter a test bankroll amount
    bankroll_input.send_keys("1000")
    
    # Find the submit button and click it
    submit_button = driver.find_element(By.CSS_SELECTOR, "button.btn-primary")
    submit_button.click()
    
    # Wait for the results to be displayed (this is a simple wait; you can implement a more sophisticated wait if needed)
    time.sleep(2)
    
    # Find the results div
    results_div = driver.find_element(By.ID, "results")
    
    # Print the results
    print("Simulation Results:")
    print(results_div.text)
    
finally:
    # Close the WebDriver
    driver.quit()
