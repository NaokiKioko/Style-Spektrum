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
    
    AllText = soup.get_text()
    AllText = AllText.replace('\n', ' ')
    AllText = AllText.replace('\t', ' ')
    
    # Get the product information
    product = GetProductInfo(AllText)

    if product is None:
        print("Failed to get product information")
        return

    # Check if the product information is correct
    if not CheckProductJson(product):
        print("Invalid product information")
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

# Main script logic
def main():
    savefolder = os.environ.get('SAVE_FOLDER', 'Images')
    continuescraping = True
    
    while continuescraping:
        print("1. Scrape product images from a webpage")
        print("2. Exit")
        print("3. Demo using (https://www.aliexpress.us/item/3256806241534476.html?src=google&gatewayAdapt=glo2usa)")
        choice = input("Enter your choice: ")

        if choice == "1":
            inputurl = input("Enter the URL of the webpage: ")
            product = ScrapeProduct(inputurl)
            print('-')
            SendMessage("product", product)
        elif choice == "2":
            continuescraping = False
        elif choice == "3":
            product = ScrapeProduct("https://www.aliexpress.us/item/3256806241534476.html?src=google&gatewayAdapt=glo2usa")
            SendMessage("product", product)
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()
