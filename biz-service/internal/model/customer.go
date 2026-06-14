package model

import "time"

type Customer struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Phone        string     `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Email        string     `gorm:"size:200" json:"email,omitempty"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	CompanyName  string     `gorm:"size:200" json:"company_name,omitempty"`
	Plan         string     `gorm:"size:20;default:'free'" json:"plan"`
	PlanExpiresAt *time.Time `json:"plan_expires_at,omitempty"`
	MaxProjects  int        `gorm:"default:1" json:"max_projects"`
	Status       string     `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type MiniappProject struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	CustomerID  uint        `gorm:"index;not null" json:"customer_id"`
	Name        string      `gorm:"size:100;not null" json:"name"`
	Description string      `gorm:"type:text" json:"description,omitempty"`
	WxAppID     string      `gorm:"size:50" json:"wx_app_id,omitempty"`
	WxAppSecret string      `gorm:"size:255" json:"-"`
	Modules     string      `gorm:"type:jsonb;default:'[]'" json:"modules"`
	BrandConfig string      `gorm:"type:jsonb;default:'{}'" json:"brand_config"`
	APIKey      string      `gorm:"uniqueIndex;size:64;not null" json:"api_key"`
	APISecret   string      `gorm:"size:64;not null" json:"-"`
	Status      string      `gorm:"size:20;default:'draft'" json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type BuildTask struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ProjectID       uint       `gorm:"index;not null" json:"project_id"`
	TriggeredBy     uint       `json:"triggered_by"`
	ModulesSnapshot string     `gorm:"type:jsonb" json:"modules_snapshot,omitempty"`
	ConfigSnapshot  string     `gorm:"type:jsonb" json:"config_snapshot,omitempty"`
	Status          string     `gorm:"size:20;default:'pending'" json:"status"`
	OutputZipURL    string     `gorm:"size:500" json:"output_zip_url,omitempty"`
	OutputMD5       string     `gorm:"size:64" json:"output_md5,omitempty"`
	ErrorLog        string     `gorm:"type:text" json:"error_log,omitempty"`
	DurationMs      int64      `json:"duration_ms,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}
