package model

import "time"

type Conversation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"index;not null" json:"tenant_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ChildID   *uint     `json:"child_id,omitempty"`
	Title     string    `gorm:"size:200" json:"title,omitempty"`
	Status    string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID uint      `gorm:"index;not null" json:"conversation_id"`
	Role           string    `gorm:"size:20;not null" json:"role"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	RetrievedChunks string   `gorm:"type:jsonb;default:'[]'" json:"retrieved_chunks,omitempty"`
	TokensUsed     int       `gorm:"default:0" json:"tokens_used"`
	CreatedAt      time.Time `json:"created_at"`
}
