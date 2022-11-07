# yt-dlp-RPC Proof of Concept

JSON-RPC 1.0 compliant RPC server for yt-dlp

### RPC Methods
- Service.Result -> details of a single process
- Service.Pending -> list of PIDs of running processes
- Service.Running -> list of details of all running processes
- Service.Kill -> kills a process by its PID

## Examples
All examples are in the **examples** folder

```js
// get all jobs
fetch('http://127.0.0.1:4444/rpc', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    'id': 0,
    'method': 'Service.Running',
    'params': []
  })
})

// start a download
fetch('http://127.0.0.1:4444/rpc', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    'id': seq,
    'method': 'Service.Exec',
    'params': [{
      'URL': this.url,
      'Params': [] //['-r', '250K'] any yt-dlp compatible args,
    }]
  })
})
```

## Notes
HTTP POST is the only availbe transport protocol at the moment, WebSocket will be implemented in the future.