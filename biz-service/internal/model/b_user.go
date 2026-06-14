package model

import "time"

type BUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null" json:"project_id"`
	OpenID    string    `gorm:"size:100;index" json:"open_id,omitempty"`
	UnionID   string    `gorm:"size:100" json:"union_id,omitempty"`
	Nickname  string    `gorm:"size:100" json:"nickname,omitempty"`
	Avatar    string    `gorm:"size:500" json:"avatar,omitempty"`
	Phone     string    `gorm:"size:20" json:"phone,omitempty"`
	Gender    int       `gorm:"default:0" json:"gender"` // 0: unknown, 1: male, 2: female
	Province  string    `gorm:"size:50" json:"province,omitempty"`
	City      string    `gorm:"size:50" json:"city,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
