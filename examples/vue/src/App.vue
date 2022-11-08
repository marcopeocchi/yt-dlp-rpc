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
      results: [],
      url: '',
      args: ''
    }
  },
  created() {
    setInterval(() => {
      fetch('http://127.0.0.1:4444/rpc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'method': 'Service.Running',
          'params': []
        })
      })
        .then(res => res.json())
        .then(data => this.results = data.result)
    }, 1000 * 1);
  },
  methods: {
    addDownload() {
      fetch('http://127.0.0.1:4444/rpc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'method': 'Service.Exec',
          'params': [{
            'URL': this.url,
            'Params': this.args.split(' ').map(a => a.trim()),
          }]
        })
      })
    },
    stopDownload(id) {
      fetch('http://127.0.0.1:4444/rpc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'method': 'Service.Kill',
          'params': [id]
        })
      })
    }
  }
}
</script>