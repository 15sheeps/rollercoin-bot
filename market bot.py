import websocket
import _thread
import time
import rel
import requests
import logging
import json
import ssl

TOKEN = ""

def purchaseItem(item):
    headers = {
        "user-agent": 'Mozilla/4.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36',
        "authorization": f'Bearer {TOKEN}',
    }

    response = requests.post('https://rollercoin.com/api/marketplace/purchase-item', headers=headers, json=item)

    print(item)
    print(response.text)

def getToken():
    headers = {
        'user-agent': 'Mozilla/4.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36',
    }

    json_data = {
        'isCaptchaRequired': False,
        'keepSigned': True,
        'mail': '123@gmail.com',
        'password': '123',
        'language': 'en',
    }

    response = requests.post('https://rollercoin.com/api/auth/email-login', headers=headers, json=json_data)

    return response.json()['data']['token']

def on_message(ws, message):
    data = json.loads(message)

    try:
      offers = data['value']['data']['list']
      total_count = data['value']['data']['totalCount']
      item_id = data['value']['item_id']

      if item_id in ["62444e4f42a0cd1b7d7bce01", "62a721f19b5a37db46bf9780", "5fd9ea94101f34db22cee8da", "644bb270648294b4642f368e"]:
        return

      item_type = data['value']['item_type']
    except:
      print("something's wrong with json:", data)
      return
    if (len(offers) > 2) and (item_type == 'miner') and (total_count >= 25):
        margin = offers[1]['price'] / offers[0]['price']
        if margin >= 1.65 and offers[1]['price'] >= 400000:
            purchaseItem({
                'itemId': item_id,
                'itemType': item_type,
                'totalCount': offers[0]['quantity'],
                'currency': 'RLT',
                'totalPrice': offers[0]['price'],
            })
            print(offers[0], offers[1], item_id, item_type)
def on_error(ws, error):
    print('error:', error)

def on_close(ws, close_status_code, close_msg):
    print("### closed ###", ws)

def on_open(ws):
    print("Opened connection")

if __name__ == "__main__":
    TOKEN = getToken()
    websocket.enableTrace(False)
    #sslopt={"cert_reqs": ssl.CERT_NONE}
    ws = websocket.WebSocketApp(f'wss://nws.rollercoin.com/?token={TOKEN}',
                              on_open=on_open,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close,
                              header = ["User-Agent: dcpgovno/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"])

    ws.run_forever(dispatcher=rel, reconnect=5)  # Set dispatcher to automatic reconnection, 5 second reconnect delay if connection closed unexpectedly
    rel.signal(2, rel.abort)  # Keyboard Interrupt
    rel.dispatch()
