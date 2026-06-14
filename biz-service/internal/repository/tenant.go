package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type TenantRepo struct{ db *gorm.DB }

func NewTenantRepo(db *gorm.DB) *TenantRepo { return &TenantRepo{db: db} }

func (r *TenantRepo) Create(t *model.Tenant) error { return r.db.Create(t).Error }

func (r *TenantRepo) FindByID(id uint) (*model.Tenant, error) {
	var t model.Tenant
	err := r.db.First(&t, id).Error
	return &t, err
}

func (r *TenantRepo) FindBySlug(slug string) (*model.Tenant, error) {
	var t model.Tenant
	err := r.db.Where("slug = ?", slug).First(&t).Error
	return &t, err
}
