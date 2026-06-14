package model

import "time"

type BProduct struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProjectID     uint      `gorm:"index;not null" json:"project_id"`
	Name          string    `gorm:"size:200;not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64   `gorm:"type:decimal(10,2)" json:"original_price,omitempty"`
	Images        string    `gorm:"type:jsonb;default:'[]'" json:"images"`
	Category      string    `gorm:"size:50" json:"category,omitempty"`
	Stock         int       `gorm:"default:0" json:"stock"`
	SalesCount    int       `gorm:"default:0" json:"sales_count"`
	AITags        string    `gorm:"type:jsonb;default:'[]'" json:"ai_tags"`
	Status        string    `gorm:"size:20;default:'on'" json:"status"` // on/off
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
