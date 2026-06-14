package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type GrowthRecord struct {
	Date       string  `json:"date"`
	Height     float64 `json:"height,omitempty"`
	Weight     float64 `json:"weight,omitempty"`
	HeadCircum float64 `json:"head_circum,omitempty"`
	Note       string  `json:"note,omitempty"`
}

type GrowthRecords []GrowthRecord

func (g GrowthRecords) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GrowthRecords) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), g)
}

type Child struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	ParentID      uint         `gorm:"index;not null" json:"parent_id"`
	TenantID      uint         `gorm:"index;not null" json:"tenant_id"`
	Name          string       `gorm:"size:100;not null" json:"name"`
	Gender        string       `gorm:"size:10" json:"gender"`
	BirthDate     time.Time    `json:"birth_date"`
	GrowthRecords GrowthRecords `gorm:"type:jsonb;default:'[]'" json:"growth_records"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
