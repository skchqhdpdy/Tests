import requests
import os
import time

st = time.time()
file_list = next((file for file in os.listdir(".") if file.endswith(".pfx")))
url = "https://www.sslshopper.com/assets/snippets/sslshopper/ajax/ajax_convert.php"
files = {'certFile': (file_list, open(file_list, 'rb'))}
body = {'certTypeFrom': 'pfx', 'certTypeTo': 'pem'}

s = requests.post(url, files=files, data=body, headers={"User-Agent": "pfx2crt.py"})
print(f"{file_list} --> {s.status_code} | www.sslshopper.com 요청 완료!")
if s.status_code == 200: s = s.content
else: print(f"{s.status_code} | 종료함!"); exit()
cert = s[s.find(b"-----BEGIN CERTIFICATE-----"):]
key = s[s.find(b"-----BEGIN PRIVATE KEY-----"):s.find(b"-----END PRIVATE KEY-----")+len(b"-----END PRIVATE KEY-----")]
with open("cert.crt", "wb") as f: f.write(cert.rstrip(b"\n"))
with open("cert.key", "wb") as f: f.write(key)
print(f"{round(time.time() - st, 2)} Sec")

changeRDPSSL = input("Replace RDP SSL? (y/n, (n)) : ").lower().startswith("y")
if changeRDPSSL:
    os.system(f'certutil -f -importpfx "{file_list}"'); time.sleep(1)
    sha1 = os.popen('powershell -Command "Get-ChildItem -Path Cert:\\LocalMachine\\My | Where-Object { $_.Subject -match \\"redstar.moe\\" } | Select-Object -ExpandProperty Thumbprint"').read().strip()
    print(sha1)
    #os.system(f'reg add "HKLM\Software\Microsoft\Terminal Server Client" /v "ServerAuthentication" /t REG_SZ /d "{sha1}" /f')
    os.system(f'reg add "HKLM\SYSTEM\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp" /v "SSLCertificateSHA1Hash" /t REG_BINARY /d "{sha1}" /f')