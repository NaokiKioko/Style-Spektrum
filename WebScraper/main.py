import os
import shutil
import json
from bs4 import BeautifulSoup
from urllib.parse import urljoin
from Requester import Requester
from ai import GetProductInfo, FilterProductPictures, TagProductImages
from sqsQueue import send_message

from dotenv import load_dotenv
load_dotenv()
requester = Requester()


def DeleteProductImages():
    # Delete productimages folder
    shutil.rmtree(requester.savefolder, ignore_errors=True)

def ScrapeSite(inputurl: str):
    response = requester.FetchHTML(inputurl)
    if response == None or response.status_code != 200:
        print("Failed to get webpage content")
        return
    soup = BeautifulSoup(response.text, 'html.parser')
    links = soup.find_all('a')
    links = set(links)
    keywords = ['product','products', 'item', ".html"]
    for link in links:
        url = link.get('href')
        if url is None:
            continue
        if any(keyword in url.lower() for keyword in keywords):
            url = urljoin(inputurl, url)
            product = ScrapeProduct(url)
            if product is None:
                print(f"Failed to scrape product information from {url}")
            else:
                send_message("StyleSpektrum", json.dumps({"Topic":"Product", "product": product}))

# Function to scrape product information from a webpage
def ScrapeProduct(url: str):
    DeleteProductImages()
    # Get the webpage content
    response = requester.FetchHTML(url)
    if response == None or response.status_code != 200:
        print("Failed to get webpage content")
        return

    # Parse the webpage content
    soup = BeautifulSoup(response.text, 'html.parser')
    AllText = extract_visible_text(soup)
    # Get the product information
    product = GetProductInfo(AllText)
    if product is None:
        print("Failed to get product information")
        return
    if product['IsCloathing'] == False:
        print("Page is not clothing")
        return
    product['URL'] = url

    # Get the product images
    image_urls = set()
    images = soup.find_all('img')

    for img in images:
        src = img.get('src')
        if not src:
            continue

        # Check conditions for relevance
        if (
            product['Title'].lower() in src.lower() or
            url.split('/')[-1].replace(".html", "").replace(".jpg", "").lower() in src.lower() or
            img.get('loading', '').lower() == "eager" or
            product['Title'].lower() in img.get('alt', '').lower()
        ):
            img_url = urljoin(url, src)
            image_urls.add(img_url)  # Using a set to avoid duplicates          

    # Filter the product images
    image_urls = list(image_urls)
    filepaths = []
    for img in image_urls:
        filepath = requester.DownloadImage(img)
        if filepath:
            filepaths.append(filepath)

    UrlsAndFolderpaths = FilterProductPictures(image_urls, filepaths)
    product['Images'] = UrlsAndFolderpaths.urls
    # FilterProductPictures(json: dict, urlList: list, folderpaths: list = None)
    product['Tags'] = TagProductImages(UrlsAndFolderpaths.filepaths)
    del product['IsCloathing']
    return product

def extract_visible_text(soup: BeautifulSoup) -> str:
    # Find all elements of interest
    content = soup.find_all(['h1', 'h2', 'h3', 'p', 'span'])  # Adjust as necessary

    # Extract and join the visible text
    visible_text = '\n'.join(element.get_text(strip=True) for element in content if element.get_text(strip=True))
    
    return visible_text

# Main script logic
def main():
    continuescraping = True
    while continuescraping:
        print("1. Scrape product")
        print("2. Scrape site")
        print("3. Exit")
        print("4. Test queue")
        choice = input("Enter your choice: ")
        if choice == "1":
            inputurl = Requester.StripDataFromURL(input("Enter the URL of the webpage: "))
            product = ScrapeProduct(inputurl)
            if product is None:
                print("Failed to scrape product information")
            else:
                send_message("StyleSpektrum", json.dumps({"Topic":"Product", "product": product}))
        elif choice == "2":
            inputurl = Requester.StripDataFromURL(input("Enter the URL of the webpage: "))
            ScrapeSite(inputurl)
        elif choice == "3":
            continuescraping = False
        elif choice == "4":
            send_message("StyleSpektrum", json.dumps({"Topic":"Test", "test": "Test message"}))
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()
