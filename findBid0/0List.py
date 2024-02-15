import os
import hashlib
import json
import time

st = time.time()

osuDir = "A:/osu!/Songs"

#beatmap_md5
def calculate_md5(filename):
    md5 = hashlib.md5()
    with open(filename, "rb") as file:
        while True:
            data = file.read(8192)
            if not data:
                break
            md5.update(data)
    return md5.hexdigest()

beatmap_info = []
try:
    for file_list in os.listdir(osuDir):
        file_list_osu = [file for file in os.listdir(f"{osuDir}/{file_list}") if file.endswith(".osu")]
        # readline_all.py
        for beatmapName in file_list_osu:
            beatmap_md5 = calculate_md5(f"{osuDir}/{file_list}/{beatmapName}")

            with open(f"{osuDir}/{file_list}/{beatmapName}", 'r', encoding="utf-8") as f:
                line = f.read()
                line = line[line.find("BeatmapID:"):]
                try:
                    bid = int(line.split("\n")[0].replace("BeatmapID:", "").replace(" ", ""))
                    bsid = int(line.split("\n")[1].replace("BeatmapSetID:", "").replace(" ", ""))
                except:
                    bid = -1
                    bsid = -1

                if bid == 0:
                    print((bid, bsid, beatmapName, beatmap_md5))
                    beatmap_info.append({"bid": bid, "bsid": bsid, "md5": beatmap_md5, "file_name": beatmapName})
except:
    print("파일나옴")
finally:
    with open("0List.json", "w") as file:
        file.write(json.dumps(beatmap_info, indent=4))
    print(f"{time.time() - st} Sec")