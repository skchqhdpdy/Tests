import aeshelper as aeshelper
import hashlib
import time
import logUtils as log

def compute_online_checksum(scoreData: list, osu_client_hash: str, storyboard_checksum: str):
    checksum_string = "chickenmcnuggets{}o15{}{}smustard{}{}uu{}{}{}{}{}{}{}Q{}{}{}{}{}{}".format(
        (int(scoreData[3])+int(scoreData[4])), scoreData[5], scoreData[6], scoreData[7], scoreData[8], 
        scoreData[0], scoreData[10], scoreData[11], scoreData[1].strip(), scoreData[9], 
        scoreData[12], scoreData[13], scoreData[14], scoreData[15], scoreData[17], scoreData[16],
        osu_client_hash, storyboard_checksum
    ).encode()
    return hashlib.md5(checksum_string).hexdigest()
    


scoreDataEnc = "PKtPuH5qZVuNSNz//P4gDy7ibwB9qwDfMYI+XcDCwZdlAhT+3mjNvSIABPRFtdDQcMsEcnzgJFqvjK6OfD79Mpgp1/aFGQwCrTbwT7LWfdYipkIPmMHWPpd1ISlzhPeif+W43eiE2ioQ20QPtHAJPhkKwkMFN+d0F9mBXtTT43A8JjQVTfA7vLCFMKAmhVs3g6lrfZHObHax4Zp7Klz5mw=="
sbk = ""
s = "Y4QQCdMoZMeiktmKchLvzfOqcrRlP9zLlNV+zOFA/aqzCSts4l1/P72dXx7AJuT9aqSbnwnODQk4C5JMwo/vAa9M7heuZhViOlACstGGKwuoznPMRwDQC3o/Pwj4OsQHUdIAf4tGPa54uPVEGmwyFjYYsgbPfsB0HrPlWzVVxLc/JZVPKCfEJBySz4nKWBQfU3Gav4GAEmIlcu4Cjvs0fXgrhY4LUa6Ex2mBTn+EeCBhyB6fMM0FjfhjIfQVdFcNCKtjfkfL9+0v3H/AK2XhgJwNoAsIKWy3YLH8F2te41Y="
iv = "J+Aeu0r7oBo1ZC4BtcpQT2sdK1NKJW8KQDA5bUQv+3Q="
aeskey = f"osu!-scoreburgr---------{20250122}"

osu_client_hash = aeshelper.decryptRinjdael(aeskey, iv, s, True) #client_hash_b64 --> decode
scoreData = aeshelper.decryptRinjdael(aeskey, iv, scoreDataEnc, True)
scoreDataSplit = scoreData.split(":")
log.debug(scoreDataSplit)
log.info(scoreData)
coc = compute_online_checksum(scoreDataSplit, osu_client_hash, sbk)
log.info(f"scoreDataSplit[2] == coc = {scoreDataSplit[2] == coc}\n{coc}")

aes_new = aeshelper.encryptRinjdael(aeskey, iv, scoreData, True)
log.info(f"scoreDataEnc == aes_new = {scoreDataEnc == aes_new}\n{aes_new}")

b64 = False
ee = aeshelper.encryptRinjdael(aeskey, iv, "1234", b64)
#log.debug(ee)
dd = aeshelper.decryptRinjdael(aeskey, iv, ee, b64)
#log.debug(dd)

time_obj = time.strptime(scoreDataSplit[16], '%y%m%d%H%M%S')
unixtime2 = int(time.mktime(time_obj) - time.timezone)
log.debug(unixtime2)

"42feb98ebb9b64188b8a5a7087d00a7d:Hi peppy :571457461768a3d71f24d8bbeb093410:2:0:0:1:0:0:660:3:True:XH:136:True:0:250218153724:20250122:50570707"
"42feb98ebb9b64188b8a5a7087d00a7d:Hi peppy :af4cef0207058cb0833fdee17d885d63:2:0:0:1:0:0:660:3:True:XH:136:True:0:250218153733:20250122:2586631"
"42feb98ebb9b64188b8a5a7087d00a7d:Hi peppy :25b7d94c1d5178067e9eb214db91dd65:2:0:0:1:0:0:660:3:True:XH:136:True:0:250218153740:20250122:8190559"
