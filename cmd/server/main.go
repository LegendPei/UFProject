package main

import (
	"UFProject/internal/config"
	"UFProject/internal/controller"
	"UFProject/internal/dao"
	"UFProject/internal/entity"
	"UFProject/internal/router"
	"UFProject/internal/service"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&entity.Graph{}, &entity.Edge{})

	graphDAO := dao.NewGraphDAO(db)
	graphService := service.NewGraphService(graphDAO)
	graphCtrl := controller.NewGraphController(graphService)

	r := router.SetupRouter(graphCtrl)
	r.Run(cfg.HTTPPort)
}
