import requests
import os
import random
from dotenv import load_dotenv
class Requester:
    savefolder = os.environ.get('SAVE_FOLDER', 'Images')

    user_agents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15',
        # Add more user agents
    ]

    def __init__(self):
        load_dotenv()
        self.proxylist = []
        filepath = os.path.join(os.path.dirname(__file__), 'proxies.txt')
        with open(filepath, 'r') as f:
            self.proxylist = f.read().splitlines()

    def GetProxy(self):
        return random.choice(self.proxylist)

    def GetHeaders(self):
        return {'User-Agent': random.choice(self.user_agents)}

    def start_requests(self, urls):
        for url in urls:
            proxy = self.GetProxy()
            html_content = self.FetchHTML(url, proxy)
            if html_content:
                print(f"Fetched HTML from {url}")

    def FetchHTML(self, url):
        attempt = False
        while (attempt == False):
            proxy = self.GetProxy()
            if len(self.proxylist) == 0:
                print("No proxies available")
                return None
            try:
                response = requests.get(url, headers=self.GetHeaders(), proxies={'http': f'http://{proxy}', 'https': f'http://{proxy}'}, timeout=100)
                response.raise_for_status()  # Raise an exception for HTTP errors
                attempt = True
                return response
            except requests.exceptions.RequestException as e:
                print(f"Error fetching {url} with proxy {proxy} | {e}\n")
                self.proxylist.remove(proxy)


    def DownloadImage(self, url: str):
        proxy = self.GetProxy()
        imageFilePaths = []
        try:
            response = requests.get(url, headers=self.GetHeaders(), proxies={'http': f'http://{proxy}', 'https': f'http://{proxy}'}, timeout=10)
            response.raise_for_status()
            image_filename = os.path.basename(url.split('?')[0])
            image_path = os.path.join(self.savefolder, image_filename)

            # Ensure folder exists
            os.makedirs(self.savefolder, exist_ok=True)

            # Save image
            with open(image_path, 'wb') as f:
                f.write(response.content)
            print(f"Image saved to {image_path}")
            imageFilePaths.append(image_path)
            return image_path
        except requests.exceptions.RequestException as e:
            print(f"Error downloading image from {url} with proxy {proxy}: {e}")
            
__all__ = ["Requester"]
