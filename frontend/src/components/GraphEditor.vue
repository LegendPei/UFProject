<template>
  <div class="editor">
    <h2>构建图</h2>
    <div class="form-group">
      <label>顶点数 (≤30):</label>
      <input
          v-model.number="nodeCount"
          type="number"
          min="1"
          max="30"
          placeholder="例如: 5"
      />
    </div>

    <h3>边列表</h3>
    <div class="table-container">
      <table>
        <thead>
        <tr>
          <th>起点</th>
          <th>终点</th>
          <th>权重 (1-99)</th>
          <th>操作</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(edge, index) in edges" :key="index">
          <td><input v-model="edge.u" placeholder="A" /></td>
          <td><input v-model="edge.v" placeholder="B" /></td>
          <td>
            <input
                v-model.number="edge.weight"
                type="number"
                min="1"
                max="99"
                placeholder="5"
            />
          </td>
          <td>
            <button @click="removeEdge(index)" class="btn-danger">删除</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <div class="button-group">
      <button @click="addEdge" class="btn-primary">+ 添加边</button>
      <button @click="saveGraph" class="btn-success">保存图</button>
      <button @click="computeMST" class="btn-warning">计算 MST</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nodeCount: 3,
      edges: [
        { u: 'A', v: 'B', weight: 5 },
        { u: 'B', v: 'C', weight: 3 }
      ]
    }
  },
  methods: {
    addEdge() {
      this.edges.push({ u: '', v: '', weight: 1 })
    },
    loadGraph(graph) {
      this.nodeCount = graph.nodeCount
      this.edges = [...graph.edges.map(e => ({
        u: e.u,
        v: e.v,
        weight: e.weight,
        isMst: e.isMst || false
      }))]
    },
    removeEdge(index) {
      this.edges.splice(index, 1)
    },
    async saveGraph() {
      try {
        const response = await fetch('/api/graphs', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            nodeCount: this.nodeCount,
            edges: this.edges
          })
        })
        if (response.ok) {
          const graph = await response.json()
          alert(`图已保存！ID: ${graph.id}`)
        } else {
          const error = await response.json()
          alert(`保存失败: ${error.error}`)
        }
      } catch (err) {
        alert('网络错误: ' + err.message)
      }
    },
    async computeMST() {
      try {
        const saveRes = await fetch('/api/graphs', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            nodeCount: this.nodeCount,
            edges: this.edges
          })
        })
        if (!saveRes.ok) {
          const error = await saveRes.json()
          throw new Error(error.error)
        }
        const graph = await saveRes.json()

        const mstRes = await fetch('/api/mst', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ graphId: graph.id })
        })
        if (!mstRes.ok) {
          const error = await mstRes.json()
          throw new Error(error.error)
        }
        const mstGraph = await mstRes.json()

        this.$emit('update-graph', mstGraph)
        alert('MST 计算完成！')
      } catch (err) {
        alert('计算失败: ' + err.message)
      }
    }
  }
}
</script>

<style scoped>
.editor {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  flex: 1;
  min-width: 300px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.table-container {
  overflow-x: auto;
  margin: 15px 0;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

th {
  background-color: #f8f9fa;
}

input {
  width: 80px;
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.button-group {
  margin-top: 15px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}

.btn-primary { background: #007bff; color: white; }
.btn-success { background: #28a745; color: white; }
.btn-warning { background: #ffc107; color: black; }
.btn-danger { background: #dc3545; color: white; }
</style>