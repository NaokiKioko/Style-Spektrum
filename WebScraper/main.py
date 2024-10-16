import os
from bs4 import BeautifulSoup
from urllib.parse import urljoin
from Requester import Requester
from ai import GetProductInfo, CheckProductJson, FilterProductPictures
from sqsQueue import SendMessage

from dotenv import load_dotenv
load_dotenv()
requester = Requester()


# Function to scrape product information from a webpage
def ScrapeProduct(url: str):
    # Get the webpage content
    response = requester.FetchHTML(url)
    if response.status_code != 200:
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

    # Get the product images
    image_urls = []
    for img in soup.find_all('img'):
        img_url = urljoin(url, img['src'])
        image_urls.append(img_url)

    # Filter the product images
    filepaths = []
    for img in image_urls:
        filepaths.append(requester.DownloadImage(img))
        
    product = FilterProductPictures(product, image_urls, filepaths)
    return product

def extract_visible_text(soup: BeautifulSoup) -> str:
    # Remove script, style, and other non-visible elements
    for element in soup(['script', 'style', 'noscript', 'meta', 'link']):
        element.extract()

    # Get all visible text
    text = soup.get_text()

    # Split into lines and remove leading/trailing spaces
    lines = (line.strip() for line in text.splitlines())

    # Break the lines further based on double spaces
    chunks = (phrase.strip() for line in lines for phrase in line.split("  "))

    # Rebuild the text by filtering out empty chunks and joining with newline
    visible_text = '\n'.join(chunk for chunk in chunks if chunk)

    return visible_text

# Main script logic
def main():
    savefolder = os.environ.get('SAVE_FOLDER', 'Images')
    continuescraping = True
    while continuescraping:
        print("1. Scrape product images from a webpage")
        print("2. Exit")
        print("3. Demo using (https://www.kohls.com/product/prd-6953168/womens-croft-barrow-34-sleeve-smocked-challis-dress.jsp)")
        choice = input("Enter your choice: ")

        if choice == "1":
            inputurl = input("Enter the URL of the webpage: ")
            product = ScrapeProduct(inputurl)
            print('-')
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
