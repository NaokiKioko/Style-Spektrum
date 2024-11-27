import os
import threading
import shutil
import json
from bs4 import BeautifulSoup
from urllib.parse import urljoin
from Requester import Requester
from ai import GetProductInfo, FilterAndTagProductPictures
from sqsQueue import send_message

from dotenv import load_dotenv
load_dotenv()

def DeleteProductImages():
    # Delete productimages folder
    shutil.rmtree(Requester.savefolder, ignore_errors=True)

def DeleteProductImages(filepaths: list):
    # Delete the image files
    for path in filepaths:
        os.remove(path)

def ExicuteScrapeProduct(url: str):
    product = ScrapeProduct(url)
    if product is None:
        print(f"Failed to scrape product information from {url}")
    else:
        send_message("StyleSpektrum", json.dumps({"Topic":"Product", "Product": product}))

def ScrapeSite(inputurl: str):
    requester = Requester()
    response = requester.FetchHTML(inputurl)
    del requester
    if response == None or response.status_code != 200:
        print("Failed to get webpage content")
        return
    soup = BeautifulSoup(response.text, 'html.parser')
    links = soup.find_all('a')
    links = set(links)
    keywords = ['product/','products/','pro/','women/','clothing/','men/', 'item/', ".html"]
    
    # Scrape each product in a new thread
    threads = []
    for link in links:
        url = link.get('href')
        if url is None:
            continue
        if any(keyword in url.lower() for keyword in keywords):
            url = urljoin(inputurl, url)
            thread = threading.Thread(target=ExicuteScrapeProduct, args=(url,))
            threads.append(thread)
            thread.start()
    for thread in threads:
        thread.join()

# Function to scrape product information from a webpage
def ScrapeProduct(url: str):
    requester = Requester()
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
            product['Name'].lower() in src.lower() or
            url.split('/')[-1].replace(".html", "").replace(".jpg", "").lower() in src.lower() or
            img.get('loading', '').lower() == "eager" or
            product['Name'].lower() in img.get('alt', '').lower() or
            product['Name'].replace(" ", "_").lower() in img.get('alt', '').lower() or
            product['Name'].replace(" ", "-").lower() in img.get('alt', '').lower()
        ):
            img_url = urljoin(url, src)
            image_urls.add(img_url)  # Using a set to avoid duplicates
        else:
            # Words in product name might be weirdly ordered in the img URL
            words = product['Name'].split()
            inCount = 0
            for word in words:
                if word.lower() in src.lower():
                    inCount += 1
            if (inCount >= len(words) - 1) and (len(words) > 2):
                img_url = urljoin(url, src)
                image_urls.add(img_url)
            elif (inCount >= len(words)) :
                img_url = urljoin(url, src)
                image_urls.add(img_url)
        
    if len(image_urls) == 0:
        print("Failed to get images")
        return None
    # Filter the product images
    image_urls = list(image_urls)
    filepaths = []
    imagesFailed = []
    for img in image_urls:
        filepath = requester.DownloadImage(img)
        if filepath:
            filepaths.append(filepath)
        else:
            imagesFailed.append(img)
    if (len(imagesFailed) > 0):
        print(f"Failed to download {len(imagesFailed)} images")
    for img in imagesFailed:
        image_urls.remove(img)
    
    UrlsAndTags = FilterAndTagProductPictures(image_urls, filepaths)
    if UrlsAndTags is None:
        print("Failed to filter and tag images")
        return None
    product['Images'] = UrlsAndTags.urls
    # FilterProductPictures(json: dict, urlList: list, folderpaths: list = None)
    product['Tags'] = UrlsAndTags.tags

    DeleteProductImages(filepaths)
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
                send_message("StyleSpektrum", json.dumps({"Topic":"Product", "Product": product}))
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
