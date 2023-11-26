import os
from mutagen.mp3 import MP3

file_list = os.listdir("songs")

def songinfo(i):
    audio = MP3(i)
    print(audio.info.sample_rate)
    return audio.info.sample_rate

def DT():
    for i in file_list:
        ffmpeg_msg = f'ffmpeg -i "songs\{i}" -af atempo=1.5 -y "DT\{i[:-4]}-DT.mp3"'
        print(f"ffmpeg_msg = {ffmpeg_msg}")
        os.system(ffmpeg_msg)

def NC():
    for i in file_list:
        pitch = songinfo(i)
        ffmpeg_msg = f'ffmpeg -i "songs\{i}" -af asetrate={pitch}*1.5 -y "NC\{i[:-4]}-NC.mp3"'
        print(f"ffmpeg_msg = {ffmpeg_msg}")
        os.system(ffmpeg_msg)

def HF():
    for i in file_list:
        ffmpeg_msg = f'ffmpeg -i "songs\{i}" -af atempo=0.75 -y "HF\{i[:-4]}-HF.mp3"'
        print(f"ffmpeg_msg = {ffmpeg_msg}")
        os.system(ffmpeg_msg)

DT()
NC()
HF()

""" for i in file_list:
    songinfo(i) """
