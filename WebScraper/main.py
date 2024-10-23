import os
import shutil
from bs4 import BeautifulSoup
from urllib.parse import urljoin
from Requester import Requester
from ai import GetProductInfo, FilterProductPictures, TagProductImages
from sqsQueue import SendMessage

from dotenv import load_dotenv
load_dotenv()
requester = Requester()


def DeleteProductImages():
    # Delete productimages folder
    shutil.rmtree(requester.savefolder, ignore_errors=True)

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
    product['url'] = url

    # Get the product images
    image_urls = []
    images = soup.find_all('img')
    for img in images:
        if product['name'] in img['src'] or url.split('/')[-1].replace(".html","").replace(".jpg", "") in img['src'] or img['loading'].lower() =="eger":
            img_url = urljoin(url, img['src'])
            if img_url in image_urls:
                continue
            image_urls.append(img_url)
        else:
            continue

    # Filter the product images
    filepaths = []
    for img in image_urls:
        filepath = requester.DownloadImage(img)
        if filepath is not None:
            filepaths.append(filepath)
    product['images'] = FilterProductPictures(product, image_urls, filepaths)
    # FilterProductPictures(json: dict, urlList: list, folderpaths: list = None)
    product['tags'] = TagProductImages(image_urls, filepaths)
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
        print("1. Scrape product images from a webpage")
        print("2. Exit")
        print("3. Demo using (https://www.kohls.com/product/prd-6953168/womens-croft-barrow-34-sleeve-smocked-challis-dress.jsp)")
        choice = input("Enter your choice: ")
        if choice == "1":
            inputurl = Requester.StripDataFromURL(input("Enter the URL of the webpage: "))
            product = ScrapeProduct(inputurl)
            if product is None:
                print("Failed to scrape product information")
            else:
                SendMessage("product", product)
        elif choice == "2":
            continuescraping = False
        elif choice == "3":
            product = ScrapeProduct("https://www.kohls.com/product/prd-6953168/womens-croft-barrow-34-sleeve-smocked-challis-dress.jsp")
            SendMessage("product", product)
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()
