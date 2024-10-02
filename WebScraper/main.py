import requests
import os
from bs4 import BeautifulSoup


def scrape_data(url):
    response = requests.get(url)
    soup = BeautifulSoup(response.text, "html.parser")
    data = soup.find_all(tag, {class_: attribute_value})
    print(data)

continue_ = True
while continue_:
    print("1. Scrape data from a website")
    print("2. Exit")
    choice = input("Enter your choice: ")
    if choice == "1":
        print("Enter the URL of the website you want to scrape data from")
        url = input()
        print("Enter the tag name you want to scrape data from")
        tag = input()
        print("Enter the class name you want to scrape data from")
        class_ = input()
        print("Enter the attribute name you want to scrape data from")
        attribute = input()
        print("Enter the attribute value you want to scrape data from")
        attribute_value = input()
        print("Enter the file name you want to save the data to")
        file_name = input()
        scrape_data(url)
    elif choice == "2":
        continue_ = False
    else:
        print("Invalid choice. Please enter a valid choice.")
