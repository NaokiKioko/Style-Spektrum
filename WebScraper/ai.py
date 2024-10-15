import google.generativeai as genai
import os

genai.configure(api_key=os.environ["API_KEY"])
model = genai.GenerativeModel("gemini-1.5-flash")

def GetProductInfo(htmldata):
    text = "With this text data could you please create a json object of the main product ignoring any 'other/similar product sections':\n"
    text += "{'title': 'Product Title', 'price': 'Product Price', 'description': 'Product Description', 'rating': 'Product Rating', 'images': ['Image URL 1', 'Image URL 2', ...], isProduct: 'bool'}\n"
    text += "The html data and text is:\n"
    text += htmldata
    response = model.generate_content(text)
    
    if not response['isProduct']:
        return None
    return response

def CheckProductHson(json):
    if 'title' in json and 'price' in json and 'description' in json and 'rating' in json and 'images' in json:
        return True
    return False


# Make Both functions available to the main.py file
__all__ = ["GetProductInfo", "CheckProductHson"]
