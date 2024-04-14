import os
import time

st = time.time()

def get_dir_size(path='.'):
    total = 0
    stack = [path]
    while stack:
        current_path = stack.pop()
        for entry in os.scandir(current_path):
            if entry.is_file():
                total += entry.stat().st_size
            elif entry.is_dir():
                stack.append(entry.path)
    if total < 1024:
        return f"{total} Byte"
    if total / (1024 ** 1) < 1024:
        return f"{round(total / (1024 ** 1), 2)} KB"
    if total / (1024 ** 2) < 1024:
        return f"{round(total / (1024 ** 2), 2)} MB"
    if total / (1024 ** 3) < 1024:
        return f"{round(total / (1024 ** 3), 2)} GB"
    if total / (1024 ** 4) < 1024:
        return f"{round(total / (1024 ** 4), 2)} TB"

print(get_dir_size("data"))

print(f"{time.time() - st} Sec")
