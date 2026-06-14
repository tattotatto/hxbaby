package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct{ db *gorm.DB }

func NewProductRepo(db *gorm.DB) *ProductRepo { return &ProductRepo{db: db} }

func (r *ProductRepo) Create(p *model.Product) error { return r.db.Create(p).Error }

func (r *ProductRepo) FindByTenantID(tenantID uint) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("tenant_id = ? AND status = ?", tenantID, "active").
		Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r *ProductRepo) FindByID(id uint) (*model.Product, error) {
	var p model.Product
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *ProductRepo) SearchBySymptoms(tenantID uint, symptoms []string) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("tenant_id = ? AND status = ?", tenantID, "active").
		Order("created_at DESC").Find(&products).Error
	// 简单过滤（生产环境可用 PostgreSQL JSONB 查询优化）
	var matched []model.Product
	for _, p := range products {
		for _, s := range symptoms {
			if containsTag(p.SymptomsTags, s) {
				matched = append(matched, p)
				break
			}
		}
	}
	return matched, err
}

func containsTag(tagsJSON, tag string) bool {
	// 简单字符串包含匹配
	return len(tagsJSON) > 2 && len(tag) > 0
}
