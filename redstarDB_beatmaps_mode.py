import requests
from dbConnent import db
import config
import time

st = time.time()
conf = config.config("config.ini")
isLog = conf.config["server"]["isLog"]
OSU_APIKEYS = eval(conf.config["osu"]["Bancho_Apikeys"])
dbR = db("redstar")

def createLog(msg):
    with open("C:/Users/Administrator/Downloads/beatmaps폴더재정렬.txt", "a", encoding="utf-8") as f: f.write(msg)
    print(msg)

data = dbR.fetch("SELECT id, rankedby, beatmap_id, beatmapset_id, beatmap_md5, song_name, mode, latest_update FROM beatmaps WHERE id >= 5390")
total = len(data)
createLog(f"{total} | DB SELECT {round(time.time() - st, 2)}s\n\n")

for i, d in enumerate(data):
    try:
        r = requests.get(f"https://osu.ppy.sh/api/get_beatmaps?k={OSU_APIKEYS[0]}&h={d['beatmap_md5']}").json()[0]
        msg = f"{i+1}/{total}({round((i+1)/total*100, 2)}%) | mode = {r['mode']} | {d['id']} | {d['beatmap_id']} | {d['song_name']} | {round(time.time() - st, 2)}s\n"
        dbR.execute(f"UPDATE beatmaps SET mode = {r['mode']} WHERE id = {d['id']}").commit()
    except: msg = f"{i+1}/{total}({round((i+1)/total*100, 2)}%) |                               ERROR! | {d['id']} | {d['beatmap_id']} | {d['song_name']} | {round(time.time() - st, 2)}s\n"
    finally: createLog(msg)

createLog(f"\nDone {round(time.time() - st, 2)}s")