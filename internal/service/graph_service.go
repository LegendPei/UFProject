package service

import (
	"UFProject/internal/dao"
	"UFProject/internal/entity"
	"sort"
)

type GraphServiceInterface interface {
	SaveGraph(nodeCount int, edges []entity.Edge) (*entity.Graph, error)
	ComputeMST(graphID uint) (*entity.Graph, error)
	ListGraphs() ([]entity.Graph, error)
	FindByID(id uint) (*entity.Graph, error)
}

var _ GraphServiceInterface = (*GraphService)(nil)

type GraphService struct {
	graphDAO dao.GraphDAOInterface
}

func NewGraphService(dao dao.GraphDAOInterface) *GraphService {
	return &GraphService{graphDAO: dao}
}

// SaveGraph 保存用户提交的原始图（不计算 MST）
// 输入：顶点数 + 边列表
// 输出：保存后的图（含自增 ID）
func (s *GraphService) SaveGraph(nodeCount int, edges []entity.Edge) (*entity.Graph, error) {
	graph := &entity.Graph{
		NodeCount: nodeCount,
		Edges:     edges,
		TotalCost: 0,
	}
	err := s.graphDAO.Save(graph)
	return graph, err
}

// ComputeMST 对指定 ID 的图运行 Kruskal 算法
// 返回计算后的图（边已标记 IsMST，TotalCost 已更新）
func (s *GraphService) ComputeMST(graphID uint) (*entity.Graph, error) {
	// 从数据库加载原始图
	graph, err := s.graphDAO.FindByID(graphID)
	if err != nil {
		return nil, err
	}

	// 执行 Kruskal 算法
	mstEdges, totalCost := s.runKruskal(graph.Edges, graph.NodeCount)

	// 更新边的 IsMST 标记
	edgeMap := make(map[string]*entity.Edge)
	for i := range graph.Edges {
		key := edgeKey(graph.Edges[i].U, graph.Edges[i].V)
		edgeMap[key] = &graph.Edges[i]
	}

	// 重置所有边的 IsMST 为 false
	for i := range graph.Edges {
		graph.Edges[i].IsMST = false
	}

	// 标记 MST 中的边
	for _, e := range mstEdges {
		key := edgeKey(e.U, e.V)
		if edge, exists := edgeMap[key]; exists {
			edge.IsMST = true
		}
	}

	// 更新总代价
	graph.TotalCost = totalCost

	// 保存结果回数据库（更新 IsMST 和 TotalCost）
	err = s.graphDAO.Save(graph)

	return graph, err
}

// edgeKey 生成边的唯一键
func edgeKey(u, v string) string {
	if u > v {
		u, v = v, u
	}
	return u + "-" + v
}

// Kruskal

// kruskalEdge 用于排序（避免修改原始 entity.Edge）
type kruskalEdge struct {
	u, v   string
	weight int
}

// runKruskal 执行 Kruskal 算法
// 输入：原始边列表，顶点数
// 输出：MST 中的边列表，总权重
func (s *GraphService) runKruskal(edges []entity.Edge, nodeCount int) ([]entity.Edge, int) {
	if nodeCount <= 1 {
		return []entity.Edge{}, 0
	}

	// 转换边格式（便于排序）
	kruskalEdges := make([]kruskalEdge, len(edges))
	for i, e := range edges {
		kruskalEdges[i] = kruskalEdge{u: e.U, v: e.V, weight: e.Weight}
	}

	// 按权重升序排序
	sort.Slice(kruskalEdges, func(i, j int) bool {
		return kruskalEdges[i].weight < kruskalEdges[j].weight
	})

	// 初始化并查集
	uf := newUnionFind(nodeCount)

	// 为顶点分配唯一 ID（因为并查集用整数）
	vertexToID := make(map[string]int)
	idCounter := 0
	for _, e := range kruskalEdges {
		if _, exists := vertexToID[e.u]; !exists {
			vertexToID[e.u] = idCounter
			idCounter++
		}
		if _, exists := vertexToID[e.v]; !exists {
			vertexToID[e.v] = idCounter
			idCounter++
		}
	}

	// Kruskal 主循环
	var mstEdges []entity.Edge
	totalCost := 0
	edgesAdded := 0

	for _, e := range kruskalEdges {
		if edgesAdded >= nodeCount-1 {
			break
		}

		uID := vertexToID[e.u]
		vID := vertexToID[e.v]

		if uf.find(uID) != uf.find(vID) {
			uf.union(uID, vID)
			mstEdges = append(mstEdges, entity.Edge{U: e.u, V: e.v, Weight: e.weight})
			totalCost += e.weight
			edgesAdded++
		}
	}

	return mstEdges, totalCost
}

// 并查集

type unionFind struct {
	parent []int
	rank   []int
}

func newUnionFind(n int) *unionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 0
	}
	return &unionFind{parent: parent, rank: rank}
}

func (uf *unionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *unionFind) union(x, y int) {
	rootX := uf.find(x)
	rootY := uf.find(y)
	if rootX != rootY {
		if uf.rank[rootX] < uf.rank[rootY] {
			uf.parent[rootX] = rootY
		} else if uf.rank[rootX] > uf.rank[rootY] {
			uf.parent[rootY] = rootX
		} else {
			uf.parent[rootY] = rootX
			uf.rank[rootX]++
		}
	}
}

// ListGraphs 获取所有图
func (s *GraphService) ListGraphs() ([]entity.Graph, error) {
	return s.graphDAO.FindAll()
}

// FindByID 根据 ID 查询图
func (s *GraphService) FindByID(id uint) (*entity.Graph, error) {
	return s.graphDAO.FindByID(id)
}
