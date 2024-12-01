from pydantic import BaseModel
import os
import json
import base64
from dotenv import load_dotenv
import re
from openai import OpenAI


load_dotenv("../.env")
# gets the api key from the environment variables
apiKey = os.environ.get('GEMINI_API_KEY')
client = OpenAI()
client.api_key = os.getenv("OPENAI_API_KEY")
# setx OPENAI_API_KEY "..."
GPTMODEL = "gpt-4o-mini"

# Function to encode the image
def encode_image(image_path):
    with open(image_path, "rb") as image_file:
        return base64.b64encode(image_file.read()).decode('utf-8')

class Product(BaseModel):
    Name: str 
    Price: float
    Description: str 
    Rating : float
    IsCloathing: bool

def GetProductInfo(htmldata: str)-> json:
    content = []  

    prompt = "You are the python def GetProductInfo(htmldata: str)->json.\n"
    prompt += "Respond with only the data of one clothing product including Name, Price, Rating. and Description if provided\n"
    prompt += "If their is not one obvious product, or the product isn't clothing set IsCloathing to false\n"
    prompt += "Using this data to fill it in:\n"
    prompt += htmldata
    
    content.append({"type": "text", "text": prompt})
    # remove tabs and newlines
    try:
        # Make the API call
        completion = client.beta.chat.completions.parse(
            model=GPTMODEL,
            messages=[{"role": "user", "content": content}],
            response_format=Product
        )
        event = completion.choices[0].message.parsed
        if event.IsCloathing == False:
            return None
        return event.model_dump()
    except Exception as e:
        print(f"Error during API call or response parsing: {e}")
    return None

def CheckProductJson(json: json):
    # check that the jsaon is acctually json
    if not isinstance(json, dict):
        return False
    # check that the json object has the correct attributes
    if not 'Name' in json:
        return False
    if not 'Price' in json:
        return False
    if not 'Description' in json:
        return False
    if not 'Rating' in json:
        return False
    if not 'IsCloathing' in json:
        return False
    return True

class UrlsAndTags(BaseModel):
    urls: list[str]
    tags: list[str]

def FilterAndTagProductPictures(url_list: list, folder_paths: list) -> UrlsAndTags:
    """
    Filters product pictures using an AI model by identifying the most frequently occurring clothing product.
    """
    content = []
    # Construct the prompt for URLs
    prompt = (
        "You are given multiple images and an array of their matching URLs. "
        "Return an array with the URLs of the clothing product that appears most frequently in the images excluding any exact duplicates or non clothing images.\n\n"
        f"URLs:\n{url_list}"
    )
    content.append({"type": "text", "text": prompt})
    # Tag image prompt
    prompt = "Can you then tag those images with the fassion styles this clothing product fits into.\m"
    prompt +="I'm asking for major categories. Do not include tags like \"short leave\", \"brown\", \"men\", \"shirt\", etc."
    content.append({"type": "text", "text": prompt})
    
    # Add encoded images to the content
    for image_path in folder_paths:
        try:
            base64_image = encode_image(image_path)  # Ensure this function is implemented
            content.append({
                "type": "image_url",
                "image_url": {
                    "url": f"data:image/jpeg;base64,{base64_image}"
                }
            })
        except Exception as e:
            print(f"Error encoding image {image_path}: {e}")
            continue

    try:
        # Make the API call
        completion = client.beta.chat.completions.parse(
            model=GPTMODEL,
            messages=[{"role": "user", "content": content}],
            response_format=UrlsAndTags
        )
        event = completion.choices[0].message.parsed
        return event

    except Exception as e:
        print(f"Error during API call or response parsing: {e}")
    return None



# Make Both functions available to the main.py file
__all__ = ["GetProductInfo", "FilterAndTagProductPictures"]
