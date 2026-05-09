import os
import time

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

BASE_URL = os.environ.get('CRAPS_URL', 'http://127.0.0.1:5000')

# Selenium 4.6+ ships with Selenium Manager — no manual chromedriver path needed.
driver = webdriver.Chrome()

try:
    driver.get(BASE_URL)

    bankroll_select = driver.find_element(By.ID, "bankroll")
    bankroll_select.click()

    submit_button = driver.find_element(By.CSS_SELECTOR, "button.btn-primary[type='submit']")
    submit_button.click()

    # Wait for the game area to become visible after starting
    WebDriverWait(driver, 5).until(
        EC.visibility_of_element_located((By.ID, "game-area"))
    )

    results_div = driver.find_element(By.ID, "results")
    print("Simulation Results:")
    print(results_div.text)

finally:
    driver.quit()
