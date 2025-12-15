package service

import (
	"UFProject/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GraphDAOInterface interface {
	Save(graph *entity.Graph) error
	FindByID(id uint) (*entity.Graph, error)
	FindAll() ([]entity.Graph, error)
	DeleteByID(id uint) error
}

type MockGraphDAO struct {
	mock.Mock
}

func (m *MockGraphDAO) Save(graph *entity.Graph) error {
	args := m.Called(graph)
	return args.Error(0)
}

func (m *MockGraphDAO) FindByID(id uint) (*entity.Graph, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Graph), args.Error(1)
}

func (m *MockGraphDAO) FindAll() ([]entity.Graph, error) {
	args := m.Called()
	return args.Get(0).([]entity.Graph), args.Error(1)
}

func (m *MockGraphDAO) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGraphService_SaveGraph(t *testing.T) {
	mockDAO := new(MockGraphDAO)
	service := NewGraphService(mockDAO)

	edges := []entity.Edge{
		{U: "A", V: "B", Weight: 5},
		{U: "B", V: "C", Weight: 3},
	}

	mockDAO.On("Save", mock.MatchedBy(func(g *entity.Graph) bool {
		return g.NodeCount == 3 && len(g.Edges) == 2
	})).Return(nil)

	graph, err := service.SaveGraph(3, edges)

	assert.NoError(t, err)
	assert.Equal(t, 3, graph.NodeCount)
	assert.Equal(t, 2, len(graph.Edges))
	mockDAO.AssertExpectations(t)
}

func TestGraphService_ComputeMST(t *testing.T) {
	mockDAO := new(MockGraphDAO)
	service := NewGraphService(mockDAO)

	originalGraph := &entity.Graph{
		ID:        1,
		NodeCount: 3,
		Edges: []entity.Edge{
			{U: "A", V: "B", Weight: 5, IsMST: false},
			{U: "B", V: "C", Weight: 3, IsMST: false},
			{U: "A", V: "C", Weight: 1, IsMST: false},
		},
	}

	mockDAO.On("FindByID", uint(1)).Return(originalGraph, nil)

	mockDAO.On("Save", mock.MatchedBy(func(g *entity.Graph) bool {
		mstEdges := 0
		total := 0
		for _, e := range g.Edges {
			if e.IsMST {
				mstEdges++
				total += e.Weight
			}
		}
		return mstEdges == 2 && total == 4
	})).Return(nil)

	result, err := service.ComputeMST(1)

	assert.NoError(t, err)
	assert.Equal(t, 4, result.TotalCost)

	mstMap := make(map[string]bool)
	for _, e := range result.Edges {
		key := e.U + "-" + e.V
		if e.IsMST {
			mstMap[key] = true
		}
	}
	assert.True(t, mstMap["A-C"] || mstMap["C-A"])
	assert.True(t, mstMap["B-C"] || mstMap["C-B"])
	assert.False(t, mstMap["A-B"] || mstMap["B-A"])

	mockDAO.AssertExpectations(t)
}
