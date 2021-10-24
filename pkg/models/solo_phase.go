package models

import "time"

type SoloPhase struct {
	ID            uint         `json:"id" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	UserID        uint         `json:"user_id" gorm:"type:int(11);column:user_id"`
	Position      string       `json:"position" gorm:"type:VARCHAR(12);column:position"`
	ValidUntil    time.Time    `json:"valid_until" gorm:"type:timestamp;column:valid_until"`
	Extensions    uint8        `json:"extensions" gorm:"type:tinyint(4);column:extensions"`
	Expired       bool         `json:"expired" gorm:"type:tinyint(1);column:expired"`
	SubdivisionID uint         `json:"-" gorm:"type:int(11);column:subdivision_id"`
	Subdivision   *Subdivision `json:"subdivision,omitempty"`
	CreatedAt     time.Time    `json:"-" gorm:"type:timestamp;column:created_at"`
	UpdatedAt     time.Time    `json:"-" gorm:"type:timestamp;column:updated_at"`
}

func (solo SoloPhase) TableName() string {
	return "atc_solo_phases"
}
