import os
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin
import urllib.request

# Function to download the product image
def download_image(image_url, folder_path):
    try:
        image_name = os.path.basename(image_url)
        image_path = os.path.join(folder_path, image_name)

        # Download and save the image
        urllib.request.urlretrieve(image_url, image_path)
        print(f"Downloaded: {image_url}")
    except Exception as e:
        print(f"Failed to download {image_url}: {e}")

# Function to scrape product images from a webpage
def scrape_product_images(url, folder_path):
    # Create folder if it doesn't exist
    if not os.path.exists(folder_path):
        os.makedirs(folder_path)

    # Send a request to fetch the content of the webpage
    response = requests.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')

    # Find all potential image tags (like <img>, <source>, etc.)
    img_tags = soup.find_all('img')

    # Loop through the image tags and filter based on product-like attributes
    for img in img_tags:
        img_url = img.get('src')

        # Some websites might use data-src for lazy-loaded images
        if not img_url:
            img_url = img.get('data-src') or img.get('data-lazy')

        # Convert relative URL to absolute URL
        img_url = urljoin(url, img_url)

        # Download the product image
        download_image(img_url, folder_path)

# Example usage: Scraping from a product page

save_folder = "product_images"
continue_scraping = True
while continue_scraping:
    print("1. Scrape product images from a webpage")
    print("2. Exit")
    choice = input("Enter your choice: ")
    if choice == "1":
        input_url = input("Enter the URL of the webpage: ")
        scrape_product_images(input_url, save_folder)