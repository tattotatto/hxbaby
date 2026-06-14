package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type ChildRepo struct{ db *gorm.DB }

func NewChildRepo(db *gorm.DB) *ChildRepo { return &ChildRepo{db: db} }

func (r *ChildRepo) Create(c *model.Child) error { return r.db.Create(c).Error }

func (r *ChildRepo) FindByParentID(parentID uint) ([]model.Child, error) {
	var children []model.Child
	err := r.db.Where("parent_id = ?", parentID).Order("created_at DESC").Find(&children).Error
	return children, err
}

func (r *ChildRepo) FindByID(id uint) (*model.Child, error) {
	var c model.Child
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *ChildRepo) Update(c *model.Child) error { return r.db.Save(c).Error }
