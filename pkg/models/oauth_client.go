package models

import "time"

type OauthClient struct {
	ID                   uint      `json:"-" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	UserID               int       `json:"user_id" gorm:"primaryKey;type:bigint(20);unsigned;column:user_id"`
	Name                 string    `json:"name" gorm:"type:varchar(255);column:name"`
	Secret               string    `json:"-" gorm:"type:varchar(100);column:secret"`
	Redirect             string    `json:"redirect" gorm:"type:varchar(200);column:redirect"`
	Revoked              bool      `json:"revoked" gorm:"type:boolean;column:revoked"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

func (client OauthClient) IsValidRedirectURI(uri string) bool {
	return client.Redirect == uri
}
