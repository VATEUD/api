package models

import "time"

type OauthClient struct {
	ID              uint      `json:"id" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	UserID          int       `json:"user_id" gorm:"primaryKey;type:bigint(20);unsigned;column:user_id"`
	Name string `json:"name" gorm:"type:varchar(255);column:name"`
	Secret string `json:"-" gorm:"type:varchar(100);column:secret"`
	
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}
