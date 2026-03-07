package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Blob struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Bucket        string         `gorm:"type:text;not null" json:"bucket"`
	Filename      string         `gorm:"type:text;not null" json:"filename"`
	Mime          string         `gorm:"type:text;not null" json:"mime"`
	Size          int64          `gorm:"type:bigint;not null" json:"size"`
	Hash          string         `gorm:"type:text;not null" json:"hash"`
	Path          string         `gorm:"type:text;not null" json:"path"`
	Public        *bool          `gorm:"type:boolean" json:"public,omitempty"`
	DownloadCount *int           `gorm:"type:int" json:"download_count,omitempty"`
	Metadata      datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
	ExpiresAt     *time.Time     `gorm:"type:timestamptz" json:"expires_at,omitempty"`
	DeletedAt     *time.Time     `gorm:"type:timestamptz" json:"deleted_at,omitempty"`
}

func (b *Blob) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
