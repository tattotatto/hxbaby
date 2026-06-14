package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(p *model.MiniappProject) error {
	return r.db.Create(p).Error
}

func (r *ProjectRepo) FindByCustomerID(customerID uint) ([]model.MiniappProject, error) {
	var projects []model.MiniappProject
	err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&projects).Error
	return projects, err
}

func (r *ProjectRepo) FindByID(id uint) (*model.MiniappProject, error) {
	var p model.MiniappProject
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) CountByCustomerID(customerID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.MiniappProject{}).Where("customer_id = ?", customerID).Count(&count).Error
	return count, err
}

func (r *ProjectRepo) Update(p *model.MiniappProject) error {
	return r.db.Save(p).Error
}
