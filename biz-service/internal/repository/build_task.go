package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type BuildTaskRepo struct {
	db *gorm.DB
}

func NewBuildTaskRepo(db *gorm.DB) *BuildTaskRepo {
	return &BuildTaskRepo{db: db}
}

func (r *BuildTaskRepo) Create(task *model.BuildTask) error {
	return r.db.Create(task).Error
}

func (r *BuildTaskRepo) FindByID(id uint) (*model.BuildTask, error) {
	var task model.BuildTask
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *BuildTaskRepo) FindByProjectID(projectID uint) ([]model.BuildTask, error) {
	var tasks []model.BuildTask
	err := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (r *BuildTaskRepo) Update(task *model.BuildTask) error {
	return r.db.Save(task).Error
}
