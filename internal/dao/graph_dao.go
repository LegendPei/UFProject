package dao

import (
	"UFProject/internal/entity"

	"gorm.io/gorm"
)

type GraphDAO struct {
	db *gorm.DB
}

func NewGraphDAO(db *gorm.DB) *GraphDAO {
	return &GraphDAO{db: db}
}

type GraphDAOInterface interface {
	Save(graph *entity.Graph) error
	FindByID(id uint) (*entity.Graph, error)
	FindAll() ([]entity.Graph, error)
	DeleteByID(id uint) error
}

var _ GraphDAOInterface = (*GraphDAO)(nil)

// Save 保存一个图
func (d *GraphDAO) Save(graph *entity.Graph) error {
	return d.db.Save(graph).Error
}

// FindByID 根据id查询图
func (d *GraphDAO) FindByID(ID uint) (*entity.Graph, error) {
	var graph entity.Graph
	err := d.db.Preload("Edges").First(&graph, ID).Error
	return &graph, err
}

// FindAll 查询所有图
func (d *GraphDAO) FindAll() ([]entity.Graph, error) {
	var graphs []entity.Graph
	err := d.db.Preload("Edges").Order("id DESC").Find(&graphs).Error
	return graphs, err
}

// DeleteByID 根据id删除图
func (d *GraphDAO) DeleteByID(id uint) error {
	if err := d.db.Where("graph_id = ?", id).Delete(&entity.Edge{}).Error; err != nil {
		return err
	}

	return d.db.Delete(&entity.Graph{}, id).Error
}
