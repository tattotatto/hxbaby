package model

import "time"

type BOrder struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   uint       `gorm:"index;not null" json:"project_id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	OrderNo     string     `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	Items       string     `gorm:"type:jsonb;not null" json:"items"`
	TotalAmount float64    `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      string     `gorm:"size:30;default:'pending'" json:"status"` // pending/paid/shipped/completed/cancelled
	Address     string     `gorm:"type:jsonb" json:"address,omitempty"`
	Remark      string     `gorm:"size:500" json:"remark,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
