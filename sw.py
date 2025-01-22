import os
import traceback
import tkinter as tk
from tkinter import messagebox

isPrivateServer = {"hosts": False, "cert": False}
hostInfo = "\n".join([
    f"175.120.229.31 akatsuki.gg",
    f"175.120.229.31 a.akatsuki.gg",
    f"175.120.229.31 admin.akatsuki.gg",
    f"175.120.229.31 assets.akatsuki.gg",
    f"175.120.229.31 b.akatsuki.gg",
    f"175.120.229.31 c.akatsuki.gg",
    f"175.120.229.31 c1.akatsuki.gg",
    f"175.120.229.31 c2.akatsuki.gg",
    f"175.120.229.31 c3.akatsuki.gg",
    f"175.120.229.31 c4.akatsuki.gg",
    f"175.120.229.31 c5.akatsuki.gg",
    f"175.120.229.31 c6.akatsuki.gg",
    f"175.120.229.31 ce.akatsuki.gg",
    f"175.120.229.31 osu.akatsuki.gg",
    f"175.120.229.31 old.akatsuki.gg",
    f"175.120.229.31 new.akatsuki.gg",
])

def get_domains():
    """
    인증서에서 지원되는 도메인(DNS Name)을 추출.
    """
    try: return [i.split("=")[1] for i in os.popen(f"certutil -user cert.crt").read().rstrip().splitlines() if "DNS Name" in i and not i.split("=")[1].startswith("*")]
    except: messagebox.showerror("ERROR!", traceback.format_exc()); return []

def on_domain_click(event):
    """
    도메인 클릭 시 실행할 작업.
    """
    try:
        # 선택된 도메인 가져오기
        selected = adl_listbox.get(adl_listbox.curselection())
        # os.system으로 작업 실행 (예: ping)
        print(f"ping {selected}")
    except: messagebox.showerror("ERROR!", traceback.format_exc())

def check_certificate():
    try: #기존 인증서 정보 확인
        r = os.popen(f"certutil -user -store root").read().rstrip()
        for i in r[r.find("* [Redstar] osu! Korea Private Server"):].splitlines():
            if "(sha1)" in i: return i.split(":")[1].strip()
    except: messagebox.showerror("ERROR!", traceback.format_exc())

def add_certificate():
    try: #인증서 추가
        messagebox.showinfo("인증서 추가 완료", os.popen(f"certutil -user -addstore root cert.crt").read().rstrip().splitlines()[-2])
        with open("C:/Windows/System32/drivers/etc/hosts", "a") as hosts: hosts.write(f"\n{hostInfo}")
        
    except: messagebox.showerror("ERROR!", traceback.format_exc())
    finally: update_button_state()

def delete_certificate():
    try: #인증서 삭제
        thumbprint = check_certificate()
        if thumbprint: messagebox.showinfo("인증서 삭제 완료", os.popen(f'certutil -user -delstore root "{thumbprint}"').read().rstrip())
        else: messagebox.showerror("결과", "인증서를 찾을 수 없습니다.")
        with open("C:/Windows/System32/drivers/etc/hosts", "r+") as hosts: 
            val = hosts.read()
            if hostInfo in val: val = val.replace(f"\n{hostInfo}", "")
            hosts.seek(0); hosts.write(val); hosts.truncate()
    except: messagebox.showerror("ERROR!", traceback.format_exc())
    finally: update_button_state()

def update_button_state(): #tkinter GUI 설정
    """현재 인증서 상태에 따라 버튼 상태 업데이트"""
    thumbprint = check_certificate()
    if thumbprint:
        delete_button.config(state=tk.NORMAL)
        add_button.config(state=tk.DISABLED)
        for domain in get_domains(): domain_listbox.insert(tk.END, domain)
    else:
        add_button.config(state=tk.NORMAL)
        delete_button.config(state=tk.DISABLED)
        domain_listbox.delete(0, tk.END)  # 기존 항목 삭제

def update_button_state(): #tkinter GUI 설정
    """현재 인증서 상태에 따라 버튼 상태 업데이트"""
    thumbprint = check_certificate()
    if thumbprint:
        delete_button.config(state=tk.NORMAL)
        add_button.config(state=tk.DISABLED)
        for domain in get_domains(): domain_listbox.insert(tk.END, domain)
    else:
        add_button.config(state=tk.NORMAL)
        delete_button.config(state=tk.DISABLED)
        domain_listbox.delete(0, tk.END)  # 기존 항목 삭제



window = tk.Tk() #tkinter 윈도우 생성
window.title("Switcher")

# 버튼 설정
add_button = tk.Button(window, text="ADD cert.crt", command=add_certificate, width=16, height=2)
delete_button = tk.Button(window, text="Remove cert.crt", command=delete_certificate, width=16, height=2)
# 버튼 배치
add_button.pack(pady=10)
delete_button.pack(pady=10)

# 도메인 목록 라벨 및 리스트박스 생성
adl = tk.Label(window, text="적용 가능한 서버 (hosts 파일 수정):")
adl_listbox = tk.Listbox(window, width=30, height=2)
adl_listbox.insert(tk.END, "akatsuki.gg --> redstar.moe")
adl_listbox.bind("<<ListboxSelect>>", on_domain_click)  # 도메인 클릭 이벤트
adl.pack(pady=5)
adl_listbox.pack(pady=5)

# 도메인 목록 라벨 및 리스트박스 생성
domain_label = tk.Label(window, text="적용 가능한 SSL 목록 (서버는 별도):")
domain_listbox = tk.Listbox(window, width=20, height=5)
domain_label.pack(pady=5)
domain_listbox.pack(pady=5)


# 초기 상태 업데이트
update_button_state()

# GUI 실행
window.mainloop()
