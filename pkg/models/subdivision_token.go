package models

import "time"

type SubdivisionToken struct {
	ID            uint   `json:"id" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	Token         string `json:"-" gorm:"type:varchar(255);column:token"`
	UserID        uint   `json:"-" gorm:"type:bigint(20);column:user_id"`
	User          User
	SubdivisionID uint `json:"-" gorm:"type:int(11);unsigned;column:subdivision_id"`
	Subdivision   Subdivision
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}
