# yt-dlp-RPC Proof of Concept

JSON-RPC 1.0 compliant RPC server for yt-dlp. Heavily inspired from Aria2 RPC.
HTTP POST and WebSockets are the available transport protocols.

## RPC Methods (overview)

| Method          | Parameters                        | Description                                                       |
|-----------------|-----------------------------------|-------------------------------------------------------------------|
| Service.Exec    | { URL: string, Params: string[] } | Starts a new process. Params is a list of yt-dlp params/arguments |
| Service.Result  | string                            | Progress of a single process                                      |
| Service.Running | -                                 | List of of all running processes progress                         |
| Service.Kill    | string                            | Kills a process by its id                                         |
| Service.KillAll | -                                 | Kills all processes                                               |

## How to run
yt-dlp path must be passed as Environmental Variable. 
```sh
# example
YT_DLP_PATH=./yt-dlp go run *.go
```
By default the server runs on port 4444
```sh
# example with different port
YT_DLP_PATH=./yt-dlp PORT=8080 go run *.go
```
### Run with Docker

```sh
docker build -t yt-dlp-rpc .
docker run --name=yt-dlp-rpc -d -p 4444:4444 -v <your directory>:/usr/src/yt-dlp-rpc/downloads yt-dlp-rpc:latest
```
## Examples
All examples are in the [examples](https://github.com/marcopeocchi/yt-dlp-rpc/tree/master/examples) folder.

```js
//HTTP POST RPC example

// get all jobs 
fetch('http://127.0.0.1:4444/rpc', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    'method': 'Service.Running',
    'params': []
  })
})

// start a download
fetch('http://127.0.0.1:4444/rpc', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    'method': 'Service.Exec',
    'params': [{
      'URL': 'https://...',
      'Params': ['-r', '500K', '--no-mtime'] // any compatible args,
    }]
  })
})
```

## Notes

Updating yt-dlp executable has to be done from the running container
```sh
docker exec -it container-name /usr/src/yt-dlp-rpc/fetch-yt-dlp.sh
```