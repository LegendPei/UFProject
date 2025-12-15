package router

import (
	"UFProject/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(graphCtrl *controller.GraphController) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/graphs", graphCtrl.SaveGraph)
		api.POST("/mst", graphCtrl.ComputeMST)
		api.GET("/graphs", graphCtrl.ListGraphs)
		api.GET("/graphs/:id", graphCtrl.GetGraph)
	}

	r.Static("/static", "./dist")
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	return r
}
