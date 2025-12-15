<template>
  <div class="history-panel">
    <h2>历史记录</h2>
    <div v-if="graphs.length === 0" class="empty">
      暂无历史记录
    </div>
    <ul v-else>
      <li v-for="graph in graphs" :key="graph.id" @click="loadGraph(graph)">
        <div class="header">
          <span>图 #{{ graph.id }}</span>
          <span>顶点: {{ graph.nodeCount }} | 边: {{ graph.edges.length }}</span>
        </div>
        <div class="stats">
          <span v-if="graph.totalCost > 0">总权重: {{ graph.totalCost }}</span>
          <span v-else>未计算 MST</span>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    return {
      graphs: []
    }
  },
  async mounted() {
    await this.loadHistory()
  },
  methods: {
    async loadHistory() {
      try {
        const res = await fetch('/api/graphs')
        this.graphs = await res.json()
      } catch (err) {
        console.error('加载历史失败:', err)
      }
    },
    loadGraph(graph) {
      // 触发父组件更新图数据
      this.$emit('load-graph', graph)
    }
  }
}
</script>

<style scoped>
.history-panel {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  flex: 1;
  min-width: 300px;
}

.empty {
  text-align: center;
  color: #6c757d;
  padding: 20px;
}

ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

li {
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: background 0.2s;
}

li:hover {
  background: #f8f9fa;
}

.header {
  display: flex;
  justify-content: space-between;
  font-weight: bold;
  margin-bottom: 5px;
}

.stats {
  font-size: 14px;
  color: #6c757d;
}
</style>