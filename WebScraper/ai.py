from pydantic import BaseModel
import google.generativeai as genai
import os
import json
import base64
from dotenv import load_dotenv
import re
from openai import OpenAI


load_dotenv("../.env")
# gets the api key from the environment variables
apiKey = os.environ.get('GEMINI_API_KEY')
genai.configure(api_key = apiKey)
model = genai.GenerativeModel("gemini-1.5-flash")
client = OpenAI()
client.api_key = os.getenv("OPENAI_API_KEY")
# setx OPENAI_API_KEY "..."
GPTMODEL = "gpt-4o-mini"

# Function to encode the image
def encode_image(image_path):
    with open(image_path, "rb") as image_file:
        return base64.b64encode(image_file.read()).decode('utf-8')


def GetProductInfo(htmldata: str)-> json:
    prompt = "You are the python def GetProductInfo(htmldata: str)->json.\n"
    prompt += "Respond with ONE and only ONE valid Json object filled with only the main products information from this data with these attributes. If their is not one obvious main product isClothing should be false:\n"
    prompt += "{'Title': string, 'Price': double, 'Description': string, 'Rating': double, IsCloathing: bool}\n"
    prompt += "Using this data to fill it in:\n"
    prompt += htmldata
    response = model.generate_content(prompt)
    # remove tabs and newlines
    text = response.text.replace("\n", "").replace("\t", "")
    if "{[" in text:
        jsonToLoad = re.search(r'[{.*}]', text).group()
        jsonToLoad = jsonToLoad.replace("[", "").replace("]", "")
    else:
        jsonToLoad = re.search(r'{.*}', text).group()
    if jsonToLoad is None:
        return None
    jsonToLoad = jsonToLoad.replace("'", "\"")
    product = json.loads(jsonToLoad)
    if product is None or product == {}:
        return None
    if not product['IsCloathing']:
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
    if not 'Title' in json:
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

class FilterEvent(BaseModel):
    urls: list[str]

class UrlsAndFolderpaths():
    urls: list[str]
    filepaths: list[str]

def FilterProductPictures(url_list: list, folder_paths: list) -> UrlsAndFolderpaths:
    """
    Filters product pictures using an AI model by identifying the most frequently occurring clothing product.
    """
    content = []
    # Construct the prompt for URLs
    prompt = (
        "You are given multiple images and an array of their matching URLs. "
        "Return an array with the URLs of the clothing product that appears most frequently in the images.\n\n"
        f"URLs:\n{url_list}"
    )
    # Add text prompt to the content
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
            response_format=FilterEvent
        )
        event = completion.choices[0].message.parsed
        
        UandF = UrlsAndFolderpaths()
        UandF.urls = event.urls
        UandF.filepaths = []
        for i, url in enumerate(url_list):  # Iterate over url_list with index
            for j, event_url in enumerate(UandF.urls):  # Iterate over UandF.urls with index
                if url == event_url:  # Compare items
                    UandF.filepaths.append(folder_paths[i])  # Append corresponding folder path
        return UandF

    except Exception as e:
        print(f"Error during API call or response parsing: {e}")
    return None

class TagEvent(BaseModel):
    tags: list[str]
def TagProductImages(folder_paths: list) -> list[str]:
    """
    tags the images with the fassion styles
    """
    content = []
    # Construct the prompt for URLs
    prompt = "can you tag these images with the fassion styles this clothing product fits into. I'm asking for major categories. Do not include tags like \"short leave\", \"brown\", \"men\", \"shirt\", etc. Please let your answer only be an single array of strings"

    # Add text prompt to the content
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
            response_format=TagEvent
        )
        event = completion.choices[0].message.parsed
        return event.tags

    except Exception as e:
        print(f"Error during API call or response parsing: {e}")
    return None




# Make Both functions available to the main.py file
__all__ = ["GetProductInfo", "FilterProductPictures", "TagProductImages"]
