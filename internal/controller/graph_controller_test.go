package controller

import (
	"UFProject/internal/entity"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type GraphServiceInterface interface {
	SaveGraph(nodeCount int, edges []entity.Edge) (*entity.Graph, error)
	ComputeMST(graphID uint) (*entity.Graph, error)
	ListGraphs() ([]entity.Graph, error)
	FindByID(id uint) (*entity.Graph, error)
}

type MockGraphService struct {
	mock.Mock
}

func (m *MockGraphService) SaveGraph(nodeCount int, edges []entity.Edge) (*entity.Graph, error) {
	args := m.Called(nodeCount, edges)
	return args.Get(0).(*entity.Graph), args.Error(1)
}

func (m *MockGraphService) ComputeMST(graphID uint) (*entity.Graph, error) {
	args := m.Called(graphID)
	return args.Get(0).(*entity.Graph), args.Error(1)
}

func (m *MockGraphService) ListGraphs() ([]entity.Graph, error) {
	args := m.Called()
	return args.Get(0).([]entity.Graph), args.Error(1)
}

func (m *MockGraphService) FindByID(id uint) (*entity.Graph, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Graph), args.Error(1)
}

func setupRouter(service *MockGraphService) *gin.Engine {
	gin.SetMode(gin.TestMode) // 静默日志
	r := gin.New()
	ctrl := NewGraphController(service)
	r.POST("/api/graphs", ctrl.SaveGraph)
	r.POST("/api/mst", ctrl.ComputeMST)
	r.GET("/api/graphs", ctrl.ListGraphs)
	r.GET("/api/graphs/:id", ctrl.GetGraph)
	return r
}

func TestGraphController_SaveGraph_Success(t *testing.T) {
	mockService := new(MockGraphService)
	router := setupRouter(mockService)

	reqBody := `{"nodeCount":3,"edges":[{"u":"A","v":"B","weight":5}]}`
	req := httptest.NewRequest("POST", "/api/graphs", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expectedGraph := &entity.Graph{
		ID:        1,
		NodeCount: 3,
		Edges: []entity.Edge{
			{ID: 1, U: "A", V: "B", Weight: 5, GraphID: 1},
		},
	}
	mockService.On("SaveGraph", 3, mock.MatchedBy(func(edges []entity.Edge) bool {
		return len(edges) == 1 && edges[0].U == "A"
	})).Return(expectedGraph, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entity.Graph
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, 3, resp.NodeCount)
	mockService.AssertExpectations(t)
}

func TestGraphController_SaveGraph_BadRequest(t *testing.T) {
	mockService := new(MockGraphService)
	router := setupRouter(mockService)

	req := httptest.NewRequest("POST", "/api/graphs", strings.NewReader(`{"nodeCount":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestGraphController_ComputeMST_Success(t *testing.T) {
	mockService := new(MockGraphService)
	router := setupRouter(mockService)

	reqBody := `{"graphId":1}`
	req := httptest.NewRequest("POST", "/api/mst", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expectedGraph := &entity.Graph{
		ID:        1,
		NodeCount: 2,
		TotalCost: 5,
		Edges: []entity.Edge{
			{U: "A", V: "B", Weight: 5, IsMST: true},
		},
	}
	mockService.On("ComputeMST", uint(1)).Return(expectedGraph, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entity.Graph
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 5, resp.TotalCost)
	assert.True(t, resp.Edges[0].IsMST)
	mockService.AssertExpectations(t)
}

func TestGraphController_GetGraph_NotFound(t *testing.T) {
	mockService := new(MockGraphService)
	router := setupRouter(mockService)

	req := httptest.NewRequest("GET", "/api/graphs/999", nil)
	w := httptest.NewRecorder()

	mockService.On("FindByID", uint(999)).Return((*entity.Graph)(nil), gorm.ErrRecordNotFound)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}
