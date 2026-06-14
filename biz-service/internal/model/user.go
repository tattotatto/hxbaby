package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"index;not null" json:"tenant_id"`
	Role      string    `gorm:"size:20;default:'parent'" json:"role"`
	Phone     string    `gorm:"size:20" json:"phone,omitempty"`
	WxOpenID  string    `gorm:"uniqueIndex;size:100" json:"wx_openid,omitempty"`
	WxUnionID string    `gorm:"size:100" json:"wx_unionid,omitempty"`
	Name      string    `gorm:"size:100" json:"name"`
	Avatar    string    `gorm:"size:500" json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
