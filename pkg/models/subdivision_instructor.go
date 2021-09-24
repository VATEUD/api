package models

import "time"

type SubdivisionInstructor struct {
	ID            uint        `json:"-" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	UserID        uint        `json:"user_id" gorm:"primaryKey;type:int(11);unsigned;unique;column:user_id"`
	Name          string      `json:"name" gorm:"type:varchar(255);column:name"`
	SubdivisionID uint        `json:"subdivision_id" gorm:"type:int(11);unsigned;column:subdivision_id"`
	Subdivision   Subdivision `json:"-"`
	CreatedAt     time.Time   `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

func (instructor SubdivisionInstructor) TableName() string {
	return "subdivision_instructors"
}
