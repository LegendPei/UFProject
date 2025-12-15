package dao

import (
	"UFProject/internal/entity"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	db.AutoMigrate(&entity.Graph{}, &entity.Edge{})
	return db
}

func TestGraphDAO_Save(t *testing.T) {
	db := setupTestDB()
	dao := NewGraphDAO(db)

	graph := &entity.Graph{
		NodeCount: 3,
		Edges: []entity.Edge{
			{U: "A", V: "B", Weight: 5},
			{U: "B", V: "C", Weight: 3},
		},
	}

	err := dao.Save(graph)
	assert.NoError(t, err)
	assert.NotZero(t, graph.ID, "Graph ID should be set")

	assert.Len(t, graph.Edges, 2)
	for _, edge := range graph.Edges {
		assert.NotZero(t, edge.ID, "Edge ID should be set")
		assert.Equal(t, graph.ID, edge.GraphID, "Edge.GraphID should match Graph.ID")
	}
}

func TestGraphDAO_FindByID(t *testing.T) {
	db := setupTestDB()
	dao := NewGraphDAO(db)

	original := &entity.Graph{
		NodeCount: 2,
		Edges: []entity.Edge{
			{U: "X", V: "Y", Weight: 10},
		},
	}
	dao.Save(original)

	found, err := dao.FindByID(original.ID)
	assert.NoError(t, err)
	assert.Equal(t, original.NodeCount, found.NodeCount)
	assert.Len(t, found.Edges, 1)
	assert.Equal(t, "X", found.Edges[0].U)
	assert.Equal(t, "Y", found.Edges[0].V)
	assert.Equal(t, 10, found.Edges[0].Weight)
}

func TestGraphDAO_FindAll(t *testing.T) {
	db := setupTestDB()
	dao := NewGraphDAO(db)

	graph1 := &entity.Graph{NodeCount: 1}
	graph2 := &entity.Graph{NodeCount: 2}
	dao.Save(graph1)
	dao.Save(graph2)

	graphs, err := dao.FindAll()
	assert.NoError(t, err)
	assert.Len(t, graphs, 2)

	assert.Equal(t, uint(2), graphs[0].ID)
	assert.Equal(t, uint(1), graphs[1].ID)
}

func TestGraphDAO_DeleteByID(t *testing.T) {
	db := setupTestDB()
	dao := NewGraphDAO(db)

	graph := &entity.Graph{
		NodeCount: 2,
		Edges: []entity.Edge{
			{U: "P", V: "Q", Weight: 7},
		},
	}
	dao.Save(graph)

	err := dao.DeleteByID(graph.ID)
	assert.NoError(t, err)

	_, err = dao.FindByID(graph.ID)
	assert.Error(t, err, "Graph should not exist after delete")

	var edges []entity.Edge
	db.Where("graph_id = ?", graph.ID).Find(&edges)
	assert.Len(t, edges, 0, "Edges should be deleted when graph is deleted")
}
