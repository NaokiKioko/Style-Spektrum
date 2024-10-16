import google.generativeai as genai
import os
import PIL.Image
import json
from dotenv import load_dotenv
import re


load_dotenv()
# gets the api key from the environment variables
apiKey = os.environ.get('GENAI_API_KEY')
genai.configure(api_key = apiKey)
model = genai.GenerativeModel("gemini-1.5-flash")

def GetProductInfo(htmldata: str):
    prompt = "Respond with a valid Json filled in with these values:\n"
    prompt += "{'title': string, 'price': int, 'description': string, 'rating': double, isCloathing: bool}\n"
    prompt += "Using this data to fill it in:\n"
    prompt += htmldata
    response = model.generate_content(prompt)
    jsonToLoad = re.search(r'{.*}', response.text).group()
    if jsonToLoad is None:
        return None
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
    # check if images are corralating to the product
    images = []
    prompt = "You are the python def FilterProductPictures(). With this product data could you please respond only with a json object with the attribute 'imagesLinks' containing the images that are related to the product information given. The order of image's matches the order of urls\n"
    prompt += "The product data is:\n"
    prompt += json['title'] + "\n"
    prompt += json['price'] + "\n"
    prompt += json['description'] + "\n"
    prompt += json['rating'] + "\n"
    prompt += "The image urls are:\n"
    for i in range(len(urlList)):
        prompt += urlList[i] + "\n"
        images.append(PIL.Image.open(folderpaths[i]))
        
    endLinks = model.generate_content(prompt,images)
    json['imagesLinks'] = endLinks
    return json


# Make Both functions available to the main.py file
__all__ = ["GetProductInfo", "CheckProductJson", "FilterProductPictures"]
