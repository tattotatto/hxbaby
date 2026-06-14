package model

import "time"

type ProductCategory struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	TenantID uint   `gorm:"index;not null" json:"tenant_id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	ParentID *uint  `json:"parent_id,omitempty"`
	SortOrder int   `gorm:"default:0" json:"sort_order"`
}

type Product struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TenantID     uint      `gorm:"index;not null" json:"tenant_id"`
	CategoryID   *uint     `json:"category_id,omitempty"`
	Name         string    `gorm:"size:200;not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description,omitempty"`
	SymptomsTags string    `gorm:"type:jsonb;default:'[]'" json:"symptoms_tags"`
	AgeRange     string    `gorm:"size:50" json:"age_range,omitempty"`
	Price        float64   `gorm:"type:decimal(10,2)" json:"price"`
	Stock        int       `gorm:"default:0" json:"stock"`
	Images       string    `gorm:"type:jsonb;default:'[]'" json:"images"`
	Status       string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
