<template>
  <div class="canvas">
    <h2>图可视化</h2>
    <div v-if="!graph" class="placeholder">
      请先构建图并计算 MST
    </div>
    <svg v-else width="500" height="400" viewBox="0 0 500 400">
      <!-- 边 -->
      <line
          v-for="(edge, index) in graph.edges"
          :key="index"
          :x1="positions[edge.u]?.x || 100"
          :y1="positions[edge.u]?.y || 100"
          :x2="positions[edge.v]?.x || 200"
          :y2="positions[edge.v]?.y || 200"
          :stroke="edge.isMst ? '#28a745' : '#6c757d'"
          :stroke-width="edge.isMst ? 3 : 1.5"
          :stroke-dasharray="!edge.isMst ? '5,5' : '0'"
      />
      <!-- 顶点 -->
      <circle
          v-for="(pos, label) in positions"
          :key="label"
          :cx="pos.x"
          :cy="pos.y"
          r="20"
          fill="#007bff"
          stroke="white"
          stroke-width="2"
      />
      <!-- 顶点标签 -->
      <text
          v-for="(pos, label) in positions"
          :key="label + '-text'"
          :x="pos.x"
          :y="pos.y + 5"
          text-anchor="middle"
          fill="white"
          font-weight="bold"
          font-size="14"
      >
        {{ label }}
      </text>
    </svg>
    <div v-if="graph" class="stats">
      <p>总权重: <strong>{{ graph.totalCost }}</strong></p>
      <p>边数: {{ graph.edges.filter(e => e.isMst).length }} / {{ graph.nodeCount - 1 }}</p>
    </div>
  </div>
</template>

<script>
export default {
  props: ['graph'],
  computed: {
    positions() {
      if (!this.graph?.edges) return {}

      const labels = new Set()
      this.graph.edges.forEach(e => {
        labels.add(e.u)
        labels.add(e.v)
      })
      const nodes = Array.from(labels)
      const positions = {}

      const centerX = 250
      const centerY = 200
      const radius = Math.min(150, 200 - nodes.length * 5)

      nodes.forEach((label, i) => {
        const angle = (i * 2 * Math.PI) / nodes.length - Math.PI / 2
        positions[label] = {
          x: centerX + radius * Math.cos(angle),
          y: centerY + radius * Math.sin(angle)
        }
      })

      return positions
    }
  }
}
</script>

<style scoped>
.canvas {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  flex: 1;
  min-width: 300px;
}

.placeholder {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6c757d;
  font-style: italic;
}

.stats {
  margin-top: 15px;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
  text-align: center;
}

.stats p {
  margin: 5px 0;
}
</style>