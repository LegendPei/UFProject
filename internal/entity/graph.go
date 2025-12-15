package entity

type Graph struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// NodeCount 表示图中顶点的数量（例如 5 个城市）
	NodeCount int `json:"nodeCount"`

	// Edges 关联的边列表（一对多关系）
	Edges []Edge `gorm:"foreignKey:GraphID" json:"edges"`

	// TotalCost MST 的总权重（计算后存入）
	TotalCost int `json:"totalCost"`
}

type Edge struct {
	ID uint `gorm:"primaryKey" json:"id"`

	GraphID uint `json:"-"`

	// U 和 V 边的两个端点（例如 "A" 和 "B"，或 "City1" 和 "City2"）
	U string `json:"u"`
	V string `json:"v"`

	// Weight 边的权重（题目要求 <100 的整数）
	Weight int `json:"weight"`

	// IsMST 标记该边是否被 Kruskal 算法选中（true = 属于最小生成树）
	IsMST bool `json:"isMst"`
}

func (Graph) TableName() string {
	return "graphs"
}

func (Edge) TableName() string {
	return "edges"
}
