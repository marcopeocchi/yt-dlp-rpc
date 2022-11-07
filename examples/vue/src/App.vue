<template>
  <main>
    <div v-if="results.length === 0">
      No active downloads
    </div>
    <div v-else v-for="result in results">
      <img v-bind:src="result.thumbnail" height="200">
      <span>{{ result.url }} - {{ result.percentage }} - {{ result.speed }}</span>
      <button @click="stopDownload(result.pid)">Stop</button>
    </div>
    <input type="text" v-model="url" />
    <button @click="addDownload">Add</button>
  </main>
</template>

<script>

export default {

  data() {
    return {
      results: [],
      url: ''
    }
  },
  created() {
    const getData = () => {
      fetch('http://127.0.0.1:4444/rpc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'id': 0,
          'method': 'Service.Running',
          'params': []
        })
      })
        .then(res => res.json())
        .then(data => this.results = data.result)
    }
    getData()
    setInterval(() => {
      getData()
    }, 1000 * 1);
  },
  methods: {
    addDownload() {
      fetch('http://127.0.0.1:4444/rpc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'id': 0,
          'method': 'Service.Exec',
          'params': [{
            'URL': this.url,
            'Params': [],
          }]
        })
      })
    },
    stopDownload(pid) {
      fetch('http://127.0.0.1:4444/rpc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'id': 0,
          'method': 'Service.Kill',
          'params': [pid]
        })
      })
    }
  }
}
</script>