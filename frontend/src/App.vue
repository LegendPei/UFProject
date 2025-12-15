<template>
  <div id="app">
    <h1>最小生成树</h1>
    <div class="main-layout">
      <div class="left-column">
        <GraphEditor
            ref="graphEditor"
            @update-graph="currentGraph = $event"
        />
        <HistoryPanel @load-graph="loadGraphFromHistory" />
      </div>
      <GraphCanvas :graph="currentGraph" />
    </div>
  </div>
</template>

<script>
import GraphEditor from './components/GraphEditor.vue'
import GraphCanvas from './components/GraphCanvas.vue'
import HistoryPanel from './components/HistoryPanel.vue'

export default {
  name: 'App',
  components: {
    GraphEditor,
    GraphCanvas,
    HistoryPanel
  },
  data() {
    return {
      currentGraph: null
    }
  },
  methods: {
    loadGraphFromHistory(graph) {
      this.currentGraph = graph
      this.$refs.graphEditor.loadGraph(graph)
    }
  }
}
</script>

<style>
.main-layout {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}

.left-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
  flex: 1;
  min-width: 300px;
}

@media (max-width: 768px) {
  .main-layout {
    flex-direction: column;
  }
}
</style>