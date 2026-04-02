import urllib.request
try:
    response = urllib.request.urlopen('https://api.github.com')
    print("Network allowed! Response:", response.read()[:100])
except Exception as e:
    print("Network blocked:", e)