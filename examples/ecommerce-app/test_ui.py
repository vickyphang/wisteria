from selenium import webdriver
import unittest

class ECommerceUITestCase(unittest.TestCase):
    def setUp(self):
        self.driver = webdriver.Chrome()  # Make sure you have the ChromeDriver installed

    def test_home_page(self):
        self.driver.get("http://localhost:8080")
        self.assertIn("Welcome to XYZ Online Retail!", self.driver.page_source)

    def tearDown(self):
        self.driver.quit()

if __name__ == "__main__":
    unittest.main()