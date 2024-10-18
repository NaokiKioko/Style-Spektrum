import requests
import os
import random
from dotenv import load_dotenv
class Requester:
    load_dotenv()
    savefolder = os.environ.get('SAVE_FOLDER', 'Images')

    user_agents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15',
        # Add more user agents
    ]

    def __init__(self):
        self.LoadProxies()

    def GetProxy(self):
        if len(self.proxylist) == 0:
            return None
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
        failCount = 0
        while (attempt == False):
            if len(self.proxylist) == 0:
                print("No proxies available")
                self.LoadProxies()
                return None
            proxy = self.GetProxy()
            try:
                # limit retry count
                response = requests.get(url, headers=self.GetHeaders(), proxies={'http': f'http://{proxy}', 'https': f'http://{proxy}'}, timeout=30)
                response.raise_for_status()  # Raise an exception for HTTP errors
                attempt = True
                self.LoadProxies()
                return response
            except requests.exceptions.RequestException as e:
                failCount += 1
                print(f"\nFail: {failCount}\nError fetching {url} with proxy {proxy} | {e}\n")
                self.proxylist.remove(proxy)
                continue


    def DownloadImage(self, url: str):
        imageFilePaths = []
        failCount = 0
        attempt = False
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

                # Save image
                with open(image_path, 'wb') as f:
                    f.write(response.content)
                print(f"Image saved to {image_path}")
                imageFilePaths.append(image_path)
                self.LoadProxies()
                return image_path
            except requests.exceptions.RequestException as e:
                failCount += 1
                print(f"\nFail: {failCount}\nError fetching {url} with proxy {proxy} | {e}\n")
                self.proxylist.remove(proxy)
                continue

    def LoadProxies(self):
        self.proxylist = []
        filepath = os.path.join(os.path.dirname(__file__), 'proxies.txt')
        with open(filepath, 'r') as f:
            self.proxylist = f.read().splitlines()
            
__all__ = ["Requester"]
