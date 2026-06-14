package model

import "time"

// CmsArticle represents a content article within a mini-app project.
type CmsArticle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Content     string    `gorm:"type:text" json:"content"`
	Summary     string    `gorm:"size:500" json:"summary,omitempty"`
	CoverImage  string    `gorm:"size:500" json:"cover_image,omitempty"`
	Category    string    `gorm:"size:50" json:"category,omitempty"`
	Tags        string    `gorm:"type:jsonb;default:'[]'" json:"tags"`
	AIGenerated bool      `gorm:"default:false" json:"ai_generated"`
	AIPrompt    string    `gorm:"type:text" json:"ai_prompt,omitempty"`
	IsPublished bool      `gorm:"default:false" json:"is_published"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
