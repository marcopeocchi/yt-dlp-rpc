import requests
import json
import time

def call(url, method, args):
  data = {
    'id': 0,
    'method': method,
    'params': [args]
  }

  res = requests.post(url=url, json=data, headers={'Content-Type': 'application/json'})
  response = json.loads(res.text)
  return response

rpc = 'http://127.0.0.1:4444/rpc'

args = {
  'URL': '...', 
  'Params': []
}

while True:
  print(call(rpc, "Service.Running", None))
  time.sleep(1)