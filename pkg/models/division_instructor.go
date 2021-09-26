package models

import "time"

type DivisionInstructor struct {
	ID        uint      `json:"-" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	UserID    uint      `json:"user_id" gorm:"primaryKey;type:int(11);unsigned;unique;column:user_id"`
	Name      string    `json:"name" gorm:"type:varchar(255);column:name"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}
