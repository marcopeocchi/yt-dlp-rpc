# yt-dlp-RPC Proof of Concept

JSON-RPC 1.0 compliant RPC server for yt-dlp

### RPC Methods (overview)
- Service.Exec    -> stars a new process
- Service.Result  -> progress of a single process
- Service.Running -> list of of all running processes progress
- Service.Kill    -> kills a process by its id
- Service.KillAll -> kills all processes

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

## Examples
All examples are in the **examples** folder

```js
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
      'URL': this.url,
      'Params': [] //['-r', '250K'] any yt-dlp compatible args,
    }]
  })
})
```

## RPC Methods

| Method          | Parameters                        | Description                                                       |
|-----------------|-----------------------------------|-------------------------------------------------------------------|
| Service.Exec    | { URL: string, Params: string[] } | Starts a new process. Params is a list of yt-dlp params/arguments |
| Service.Result  | string                            | Progress of a single process                                      |
| Service.Running | -                                 | List of of all running processes progress                         |
| Service.Kill    | string                            | Kills a process by its id                                         |
| Service.KillAll | -                                 | Kills all processes                                               |