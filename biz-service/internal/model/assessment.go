package model

import "time"

type Assessment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID uint      `gorm:"index;not null" json:"conversation_id"`
	ChildID        *uint     `json:"child_id,omitempty"`
	TenantID       uint      `gorm:"index;not null" json:"tenant_id"`
	Symptoms       string    `gorm:"type:jsonb;default:'[]'" json:"symptoms"`
	AIAnalysis     string    `gorm:"type:text" json:"ai_analysis,omitempty"`
	RiskLevel      string    `gorm:"size:20" json:"risk_level,omitempty"`
	SuggestionType string    `gorm:"size:30" json:"suggestion_type,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type Recommendation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MessageID    uint      `gorm:"index;not null" json:"message_id"`
	ProductID    uint      `json:"product_id"`
	Score        float64   `json:"score"`
	Reason       string    `gorm:"type:text" json:"reason,omitempty"`
	UserFeedback string    `gorm:"size:20" json:"user_feedback,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
