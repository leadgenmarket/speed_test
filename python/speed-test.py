import json
import pingparsing
import speedtest
import time
import requests
from email.mime.text import MIMEText
from email.header    import Header
import atexit

import ssl
ssl._create_default_https_context = ssl._create_unverified_context

def testInternet(user, limit, ping, ping_limit) :
    ping_parser = pingparsing.PingParsing()
    transmitter = pingparsing.PingTransmitter()
    transmitter.destination = ping
    transmitter.count = 50
    result = transmitter.ping()
    params = ping_parser.parse(result)
    st = speedtest.Speedtest()
    st.get_best_server()
    speedDownload = st.download()/1000000
    speedUpload = st.upload()/1000000
    ping = st.results.ping
    

    packetLossPercentage = (params.packet_loss_count/50)*100

    sendRequest(str(int(packetLossPercentage)), str(int(speedDownload)), str(int(speedUpload)), user, str(ping))

def sendRequest(packet_loss_count, speedDownload, speedUpload, user, ping) :
    r = requests.get("http://mail.leadactiv.ru/saveData.php?user="+user+"&insp="+speedDownload+"&outsp="+speedUpload+"&ping="+ping+"&loss="+packet_loss_count+"&key=swGh889KyxjWyz")
    

f = open('config.json', 'r', encoding='utf-8')
config= json.load(f)

while True:
    try:
        testInternet(config["operator"], config["speed_limit"], config["ping"], config["ping_limit"])
    except Exception: 
        pass
    time.sleep(config["repeat_mins"]*60)








