package controller

import (
	"UFProject/internal/entity"
	"UFProject/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GraphController struct {
	graphService service.GraphServiceInterface
}

func NewGraphController(service service.GraphServiceInterface) *GraphController {
	return &GraphController{graphService: service}
}

// SaveGraph 处理 POST /api/graphs
// 请求体：{ "nodeCount": 5, "edges": [{"u":"A","v":"B","weight":5}, ...] }
func (c *GraphController) SaveGraph(ctx *gin.Context) {
	var req struct {
		NodeCount int           `json:"nodeCount" binding:"required,min=1,max=30"`
		Edges     []entity.Edge `json:"edges" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	graph, err := c.graphService.SaveGraph(req.NodeCount, req.Edges)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save graph"})
		return
	}

	ctx.JSON(http.StatusOK, graph)
}

// ComputeMST 处理 POST /api/mst
// 请求体：{ "graphId": 1 }
func (c *GraphController) ComputeMST(ctx *gin.Context) {
	var req struct {
		GraphID uint `json:"graphId" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	graph, err := c.graphService.ComputeMST(req.GraphID)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Graph not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute MST"})
		}
		return
	}

	ctx.JSON(http.StatusOK, graph)
}

// ListGraphs 处理 GET /api/graphs
func (c *GraphController) ListGraphs(ctx *gin.Context) {
	graphs, err := c.graphService.ListGraphs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list graphs"})
		return
	}

	ctx.JSON(http.StatusOK, graphs)
}

// GetGraph 处理 GET /api/graphs/:id
func (c *GraphController) GetGraph(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid graph ID"})
		return
	}

	graph, err := c.graphService.FindByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Graph not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get graph"})
		}
		return
	}

	ctx.JSON(http.StatusOK, graph)
}
