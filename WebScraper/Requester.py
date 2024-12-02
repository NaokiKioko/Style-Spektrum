import requests
import os
import random
from dotenv import load_dotenv
# from selenium import webdriver
# from selenium.webdriver.chrome.options import Options

class Requester:
    load_dotenv()
    savefolder = os.environ.get('SAVE_FOLDER', os.path.join(os.path.dirname(__file__), 'images'))

    user_agents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15',
        # Add more user agents
    ]
    working_proxie = None

    def __init__(self):
        self.LoadProxies()
        os.makedirs(self.savefolder, exist_ok=True)

    def GetProxy(self):
        if self.working_proxie:
            return self.working_proxie
        if len(self.proxylist) == 0:
            return None
        return random.choice(self.proxylist)

    def GetHeaders(self, ):
        return {
            'User-Agent': random.choice(self.user_agents),
            'Accept-Language': 'en-US,en;q=0.9',
            'Connection': 'keep-alive',
        }

    def start_requests(self, urls):
        for url in urls:
            proxy = self.GetProxy()
            html_content = self.FetchHTML(url, proxy)
            if html_content:
                print(f"Fetched HTML from {url}")

    def FetchHTML(self, url):
        attempt = False
        failCount = 0
        while (attempt == False):
            if len(self.proxylist) == 0:
                print("No proxies available")
                self.LoadProxies()
                return None
            proxy = self.GetProxy()
            try:
                # limit retry count
                response = requests.get(url, headers=self.GetHeaders(), proxies={'http': f'http://{proxy}', 'https': f'http://{proxy}'}, timeout=10)
                response.raise_for_status()  # Raise an exception for HTTP errors
                attempt = True
                self.working_proxie = proxy
                self.LoadProxies()
                print(f"HTML fetched from {url} Successfully")
                return response
            except requests.exceptions.RequestException as e:
                failCount += 1
                print(f"\nProxies Failed: {failCount}\nError fetching {url} with proxy {proxy} | {e}\n")
                # Set working proxy
                self.working_proxie = None
                self.proxylist.remove(proxy)
                continue

    def DownloadImage(self, url: str):
        failCount = 0
        attempt = False
        self.LoadProxies()
        while (attempt == False):
            proxy = self.GetProxy()
            try:
                response = requests.get(url, headers=self.GetHeaders(), proxies={'http': f'http://{proxy}', 'https': f'http://{proxy}'}, timeout=10)
                response.raise_for_status()
                attempt = True
                image_filename = os.path.basename(url.split('?')[0])
                image_path = os.path.join(self.savefolder, image_filename)
                # Ensure folder exists
                os.makedirs(self.savefolder, exist_ok=True)
                # Make sure the image path is unique
                image_path = self.ValidImagePath(image_path)
                # Save image
                with open(image_path, 'wb') as f:
                    f.write(response.content)
                print(f"Image saved to {image_path}")
                # Set working proxy
                self.working_proxie = proxy
                self.LoadProxies()
                return image_path
            except requests.exceptions.RequestException as e:
                failCount += 1
                # print(f"\nProxies Failed: {failCount}\nError fetching {url} with proxy {proxy} | {e}\n")
                if self.proxylist.count(proxy) > 0:
                    self.working_proxie = None
                    self.proxylist.remove(proxy)
                else:
                    return None
                continue
    
    def ValidImagePath(self, image_path: str) -> str:
        count = 0
        original_path =image_path.rsplit('.', 1)  # Preserve the original path for appending

        while os.path.isfile(image_path):
            count += 1
            
            image_path = f"{original_path[0]}_{count}.{original_path[1]}"

        return image_path
    
    def LoadProxies(self):
        self.proxylist = []
        filepath = os.path.join(os.path.dirname(__file__), 'proxies.txt')
        with open(filepath, 'r') as f:
            self.proxylist = f.read().splitlines()
            
    def StripDataFromURL(url: str):
        return url.split('?')[0]
__all__ = ["Requester"]
