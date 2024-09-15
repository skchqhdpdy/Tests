import requests
import os

url = "https://www.sslshopper.com/assets/snippets/sslshopper/ajax/ajax_convert.php"
file_list = next((file for file in os.listdir(".") if file.endswith(".pfx")))

files = {'certFile': (file_list, open(file_list, 'rb'))}
data = {'certTypeFrom': 'pfx', 'certTypeTo': 'pem'}

with open("ssl.pem", "wb") as f: f.write(requests.post(url, files=files, data=data).content)

with open("ssl.pem", "r") as s:
    s = s.read()
    key = s[s.find("-----BEGIN PRIVATE KEY-----"):s.find("-----END PRIVATE KEY-----")+len("-----END PRIVATE KEY-----")]
    with open("cert.key", "w") as f: f.write(key)

    cert = s[s.find("-----BEGIN CERTIFICATE-----"):].rstrip("\n")
    with open("cert.crt", "w") as f: f.write(cert)

os.remove("ssl.pem")   