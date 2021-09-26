package models

import (
	"time"
)

type News struct {
	ID        uint      `json:"id" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	Title     string    `json:"title" gorm:"type:varchar(255);column:title"`
	Body      string    `json:"body" gorm:"type:varchar(65535);column:body"`
	UserID    uint      `json:"user_id" gorm:"primaryKey;type:int(11);unsigned;column:user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

func (news News) TableName() string {
	return "news"
}
