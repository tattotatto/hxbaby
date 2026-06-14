package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

// CmsArticleRepo handles database operations for CMS articles.
type CmsArticleRepo struct {
	db *gorm.DB
}

// NewCmsArticleRepo creates a new CmsArticleRepo with the given DB connection.
func NewCmsArticleRepo(db *gorm.DB) *CmsArticleRepo {
	return &CmsArticleRepo{db: db}
}

// List returns paginated published articles for a project, optionally filtered by category.
func (r *CmsArticleRepo) List(projectID uint, category string, page, pageSize int) ([]model.CmsArticle, int64, error) {
	var articles []model.CmsArticle
	var total int64
	query := r.db.Where("project_id = ? AND is_published = ?", projectID, true)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query.Model(&model.CmsArticle{}).Count(&total)
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error
	return articles, total, err
}

// FindByID retrieves a single article by ID scoped to the project.
func (r *CmsArticleRepo) FindByID(id, projectID uint) (*model.CmsArticle, error) {
	var article model.CmsArticle
	err := r.db.Where("id = ? AND project_id = ?", id, projectID).First(&article).Error
	return &article, err
}

// Create inserts a new CMS article into the database.
func (r *CmsArticleRepo) Create(article *model.CmsArticle) error {
	return r.db.Create(article).Error
}

// IncrementViewCount atomically increments the view_count column for the given article.
func (r *CmsArticleRepo) IncrementViewCount(id uint) error {
	return r.db.Model(&model.CmsArticle{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
