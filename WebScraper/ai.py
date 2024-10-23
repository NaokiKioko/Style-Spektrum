from tkinter import Image
import google.generativeai as genai
import os
import json
import requests
from PIL import Image
import PIL
from dotenv import load_dotenv
import re


load_dotenv("../.env")
# gets the api key from the environment variables
apiKey = os.environ.get('GEMINI_API_KEY')
genai.configure(api_key = apiKey)
model = genai.GenerativeModel("gemini-1.5-flash")

def GetProductInfo(htmldata: str)-> json:
    prompt = "You are the python def GetProductInfo(htmldata: str)->json.\n"
    prompt += "Respond with one and only one valid Json filled with only the main products information from this data with these attributes:\n"
    prompt += "{'title': string, 'price': double, 'description': string, 'rating': double, isCloathing: bool}\n"
    prompt += "Using this data to fill it in:\n"
    prompt += htmldata
    response = model.generate_content(prompt)
    # remove tabs and newlines
    text = response.text.replace("\n", "").replace("\t", "")
    jsonToLoad = re.search(r'{.*}', text).group()
    if jsonToLoad is None:
        return None
    jasonToLoad = jsonToLoad.replace("'", "\"")
    product = json.loads(jsonToLoad)
    if not product['isCloathing']:
        return None
    # confirm that the json object is correct
    if not CheckProductJson(product):
        return None
    return product

def CheckProductJson(json: json):
    # check that the jsaon is acctually json
    if not isinstance(json, dict):
        return False
    # check that the json object has the correct attributes
    if not 'title' in json:
        return False
    if not 'price' in json:
        return False
    if not 'description' in json:
        return False
    if not 'rating' in json:
        return False
    if not 'isCloathing' in json:
        return False
    return True

def FilterProductPictures(json: json, urlList: list, folderpaths: list):
    if len(urlList) != len(folderpaths):
        raise ValueError("urlList and folderpaths must have the same length.")
    images = [Image.open(path) for path in folderpaths]

    prompt = "You are the python def FilterProductPictures(). With this product data and images, please respond only with a json object with the attribute 'imagesLinks' containing the URLs of the images that showcase the same clothing item. \n"
    prompt += f"The product data is:\n"
    prompt += f"{json['title']}\n"
    prompt += f"{json['price']}\n"
    prompt += f"{json['description']}\n"
    prompt += f"{json['rating']}\n"
    prompt += "The image URLs are:\n"
    for i, url in enumerate(urlList):
        prompt += f"{i+1}. {url}\n"
    images.append(prompt)
    response = model.generate_content(images)
    try:
        # Parse JSON response using json library
        data = json.loads(response.text)
        return data.get('imagesLinks')
    except (json.JSONDecodeError, KeyError):
        # Handle errors during parsing
        print("Error parsing generated JSON")
        return None

def TagProductImages(urlList: list, folderpaths: list):
    # check if images are corralating to the product
    images = []
    prompt = "You are the python def TagProductImages(). With these product Images please respond with a json object with the attribute 'tags' containing the fassion styles these images fit into\n"
    for i in range(len(urlList)):
        prompt += urlList[i] + "\n"
        images.append(PIL.Image.open(folderpaths[i]))
    images.append(prompt)
    response = model.generate_content(images)
    jsonToLoad = re.search(r'{.*}', response.text).group()
    if jsonToLoad is None:
        return None
    tags = json.loads(jsonToLoad)
    return tags


# Make Both functions available to the main.py file
__all__ = ["GetProductInfo", "FilterProductPictures", "TagProductImages"]
