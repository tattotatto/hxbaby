package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type ConversationRepo struct{ db *gorm.DB }

func NewConversationRepo(db *gorm.DB) *ConversationRepo { return &ConversationRepo{db: db} }

func (r *ConversationRepo) Create(c *model.Conversation) error { return r.db.Create(c).Error }

func (r *ConversationRepo) FindByTenantAndUser(tenantID, userID uint) ([]model.Conversation, error) {
	var convs []model.Conversation
	err := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Order("updated_at DESC").Find(&convs).Error
	return convs, err
}

func (r *ConversationRepo) FindByID(id uint) (*model.Conversation, error) {
	var c model.Conversation
	err := r.db.First(&c, id).Error
	return &c, err
}

type MessageRepo struct{ db *gorm.DB }

func NewMessageRepo(db *gorm.DB) *MessageRepo { return &MessageRepo{db: db} }

func (r *MessageRepo) Create(m *model.Message) error { return r.db.Create(m).Error }

func (r *MessageRepo) FindByConversationID(convID uint, limit int) ([]model.Message, error) {
	var msgs []model.Message
	err := r.db.Where("conversation_id = ?", convID).
		Order("created_at ASC").Limit(limit).Find(&msgs).Error
	return msgs, err
}
