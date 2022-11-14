import requests
import json
import time
import sys

# HTTP POST python3 RPC client example

def call(url, method, args):
  data = {
    'id': 0,
    'method': method,
    'params': [args]
  }

  res = requests.post(url=url, json=data, headers={'Content-Type': 'application/json'})
  response = json.loads(res.text)
  return response

if len(sys.argv) != 3:
  print('usage example: python3 client.py localhost:4444 https://videoresource.example')
  sys.exit(1)

addr = sys.argv[1]
url  = sys.argv[2]

rpc = f'http://{addr}/rpc'
args = {
  'URL': url, # type a video url
  'Params': []
}

call(rpc, 'Service.Exec', args)

while True:
  res = call(rpc, 'Service.Running', None)
  for r in res['result']:
    eid = r['id']
    info = r['info']
    progress = r['progress']

    url = info['url']
    percentage = progress['percentage']
    
    if percentage == '-1':
      print(f'{url} Done!')
      call(rpc, 'Service.Kill', eid)
    else:
      print(f'{id} {url} {percentage}')

  time.sleep(1)
