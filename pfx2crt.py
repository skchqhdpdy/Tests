import requests
import os

file_list = next((file for file in os.listdir(".") if file.endswith(".pfx")))
url = "https://www.sslshopper.com/assets/snippets/sslshopper/ajax/ajax_convert.php"
files = {'certFile': (file_list, open(file_list, 'rb'))}
body = {'certTypeFrom': 'pfx', 'certTypeTo': 'pem'}

s = requests.post(url, files=files, data=body, headers={"User-Agent": "pfx2crt.py"}).content
cert = s[s.find(b"-----BEGIN CERTIFICATE-----"):]
key = s[s.find(b"-----BEGIN PRIVATE KEY-----"):s.find(b"-----END PRIVATE KEY-----")+len(b"-----END PRIVATE KEY-----")]
with open("cert.crt", "wb") as f: f.write(cert.rstrip(b"\n"))
with open("cert.key", "wb") as f: f.write(key)