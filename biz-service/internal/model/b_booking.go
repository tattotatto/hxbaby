package model

import "time"

type BBooking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null" json:"project_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	SlotTime  time.Time `json:"slot_time"`
	Duration  int       `gorm:"default:30" json:"duration"` // minutes
	Name      string    `gorm:"size:50;not null" json:"name"`
	Phone     string    `gorm:"size:20;not null" json:"phone"`
	Remark    string    `gorm:"size:500" json:"remark,omitempty"`
	Status    string    `gorm:"size:20;default:'confirmed'" json:"status"` // confirmed/cancelled/completed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
