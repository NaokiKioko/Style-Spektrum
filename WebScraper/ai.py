import google.generativeai as genai
import os
import PIL.Image

# gets the api key from the environment variables
apiKey = os.environ.get('GENAI_API_KEY')
genai.configure(api_key = apiKey)
model = genai.GenerativeModel("gemini-1.5-flash")

def GetProductInfo(htmldata):
    prompt = "With this text data could you please create a json object of the main product ignoring any 'other/similar product sections'. The isProduct :\n"
    prompt += "{'title': 'Product Title', 'price': 'Product Price', 'description': 'Product Description', 'rating': 'Product Rating', isProduct: 'bool'}\n"
    prompt += "The html data and text is:\n"
    prompt += htmldata
    response = model.generate_content(prompt)
    
    if not response['isProduct']:
        return None
    return response

def CheckProductJson(json):
    if 'title' in json and 'price' in json and 'description' in json and 'rating' in json and 'images' in json:
        return True
    return False

def FilterProductPictures(json, urlList, folderpaths):
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
