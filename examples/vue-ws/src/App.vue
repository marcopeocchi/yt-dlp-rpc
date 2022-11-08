<template>
  <main>
    <div v-if="results.length === 0">
      No active downloads
    </div>
    <div v-else v-for="result in results">
      <img v-bind:src="result.thumbnail" height="200">
      <span>{{ result.url }} - {{ result.percentage }} - {{ result.speed }}</span>
      <button @click="stopDownload(result.id)">Stop</button>
    </div>
    <label>url</label>
    <input type="text" v-model="url" />
    <button @click="addDownload">Add</button>
    <div>
      <label>args</label>
      <input type="text" v-model="args" />
    </div>
  </main>
</template>

<script>

export default {

  data() {
    return {
      url: '',
      args: '',
      results: [],
      socket: new WebSocket('ws://localhost:4444/rpc-ws')
    }
  },
  created() {
    const getRunning = () => this.socket.send(JSON.stringify({
      'method': 'Service.Running',
      'params': []
    }))

    this.socket.onopen = () => getRunning

    this.socket.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.results = data.result
    }

    setInterval(() => {
      getRunning()
    }, 1250);
  },
  methods: {
    addDownload() {
      this.socket.send(JSON.stringify({
        'method': 'Service.Exec',
        'params': [{
          'URL': this.url,
          'Params': this.args.split(' ').map(a => a.trim()),
        }]
      }))
    },
    stopDownload(id) {
      this.socket.send(JSON.stringify({
        'method': 'Service.Kill',
        'params': [id]
      }))
    }
  }
}
</script>