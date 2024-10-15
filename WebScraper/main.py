import os
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin
import urllib.request
from ai import GetProductInfo, CheckProductHson


# Function to download the product image
def DownloadImage(imageurl, folderpath):
    # Create folder for images if it doesn't exist
    if not os.path.exists(folderpath):
        os.makedirs(folderpath)

    try:
        imageName = os.path.basename(imageurl)
        imagePath = os.path.join(folderpath, imageName)

        # Download and save the image
        urllib.request.urlretrieve(imageurl, imagePath)
        print(f"Downloaded: {imageurl}")
    except Exception as e:
        print(f"Failed to download {imageurl}: {e}")

def HandleImages(soup, url, folderpath):
    # Find all potential image tags (like <img>, <source>, etc.)
    imgTags = soup.find_all('img')
    imageLinks = []
    # Loop through the image tags and filter based on product-like attributes
    for img in imgTags:
        imgurl = img.get('src')

        # Some websites might use data-src for lazy-loaded images
        if not imgurl:
            imgurl = img.get('data-src') or img.get('data-lazy')

        # Convert relative URL to absolute URL
        imgurl = urljoin(url, imgurl)
        imageLinks.append(imgurl)
        
        # Download the product image
        DownloadImage(imgurl, folderpath)
    return imageLinks

def handleproductinfo(soup, imagelinks):
    body = soup.find('body')
    PageText = body.gettext()
    
    product = GetProductInfo(PageText)
    if not CheckProductHson(product):
        print("Failed to get product info from AI")
        return None
    
    product = {
        "title": product['title'],
        "price": product['price'],
        "description": product['description'],
        "rating": product['rating'],
        "images": imagelinks
    }
    return product

# Function to scrape product from a webpage
def ScrapeProduct(url, folderpath):
    # Send a request to fetch the content of the webpage
    response = requests.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')
    imagelinks = HandleImages(soup, url, folderpath)
    
    product = handleproductinfo(soup, imagelinks)
    print("This is the product we got")
    print(product)

# ----- Global Varibles -----
savefolder = "productimages"
continuescraping = True
# ----- Global Varibles -----

# -------------------- MAIN --------------------
while continuescraping:
    print("1. Scrape product images from a webpage")
    print("2. Exit")
    print("3. Demo using (https://books.toscrape.com/)")
    choice = input("Enter your choice: ")
    if choice == "1":
        inputurl = input("Enter the URL of the webpage: ")
        product = ScrapeProduct(inputurl, savefolder)
        # Put together GPT Prompt
        # ask gpt to generate tags for the images
        # Place results in msg queue
        # delete all images in the folder
    elif choice == "2":
        continuescraping = False
    elif choice == "3":
        ScrapeProduct("https://books.toscrape.com/", savefolder)
    else:
        print("Invalid choice. Please try again.")