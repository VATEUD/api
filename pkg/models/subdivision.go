package models

import "time"

type Subdivision struct {
	ID        uint      `json:"id" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	Code     string    `json:"code" gorm:"type:varchar(20);column:code"`
	Name      string    `json:"name" gorm:"type:varchar(255);column:name"`
	WebsiteURL      string    `json:"website_url" gorm:"type:varchar(255);column:website_url"`
	DiscordServerURL      string    `json:"discord_server_url" gorm:"type:varchar(100);column:discord_server_url"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}