package models

import "time"

type OauthAuthCode struct {
	ID        string `json:"-" gorm:"primaryKey;type:varchar(100);unsigned;column:id;unique"`
	ClientID  uint   `json:"client_id" gorm:"type:bigint(20);unsigned;column:client_id"`
	Client    OauthClient
	Scopes    string    `json:"-" gorm:"type:varchar(200);column:scopes"`
	UserAgent string    `json:"-" gorm:"type:varchar(200);column:user_agent"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}
