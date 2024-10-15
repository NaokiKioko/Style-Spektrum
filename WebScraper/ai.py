import google.generativeai as genai
import os
import PIL.Image
import json

# gets the api key from the environment variables
apiKey = os.environ.get('GENAI_API_KEY')
genai.configure(api_key = apiKey)
model = genai.GenerativeModel("gemini-1.5-flash")

def GetProductInfo(htmldata: str):
    prompt = "You are the python def GetProductInfo().With this text data create a json object of the main product ignoring any 'other/similar product' sections:\n"
    prompt += "{'title': 'Product Title', 'price': 'Product Price', 'description': 'Product Description', 'rating': 'Product Rating', isCloathing: 'bool'}\n"
    prompt += "The html data and text is:\n"
    prompt += htmldata
    response = model.generate_content(prompt)
    product = json.loads(response)
    if not response['isCloathing']:
        return None
    # confirm that the json object is correct
    if not CheckProductJson(response):
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
    if not 'isProduct' in json:
        return False
    return True

def FilterProductPictures(json: json, urlList: list, folderpaths: list):
    # check if images are corralating to the product
    images = []
    prompt = "You are the python def FilterProductPictures(). With this product data could you please respon only with a json object with the attribute 'imagesLinks' containing the images that are related to the product information given.\n"
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
