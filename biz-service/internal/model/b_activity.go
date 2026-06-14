package model

import "time"

type BActivity struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ProjectID      uint       `gorm:"index;not null" json:"project_id"`
	Title          string     `gorm:"size:200;not null" json:"title"`
	Description    string     `gorm:"type:text" json:"description"`
	CoverImage     string     `gorm:"size:500" json:"cover_image,omitempty"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	SignupStart    *time.Time `json:"signup_start,omitempty"`
	SignupEnd      *time.Time `json:"signup_end,omitempty"`
	MaxParticipants int       `gorm:"default:0" json:"max_participants"`
	CurrentCount   int        `gorm:"default:0" json:"current_count"`
	Location       string     `gorm:"size:300" json:"location,omitempty"`
	AICopy         string     `gorm:"type:text" json:"ai_copy,omitempty"`
	AIPoster       string     `gorm:"size:500" json:"ai_poster,omitempty"`
	Status         string     `gorm:"size:20;default:'draft'" json:"status"` // draft/published/ended
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type BActivitySignup struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ActivityID  uint       `gorm:"index;not null" json:"activity_id"`
	ProjectID   uint       `gorm:"index;not null" json:"project_id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	Name        string     `gorm:"size:50;not null" json:"name"`
	Phone       string     `gorm:"size:20;not null" json:"phone"`
	Remark      string     `gorm:"size:300" json:"remark,omitempty"`
	CheckedIn   bool       `gorm:"default:false" json:"checked_in"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
