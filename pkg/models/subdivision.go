package models

import "time"

type Subdivision struct {
	ID               uint                    `json:"-" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	Code             string                  `json:"code" gorm:"type:varchar(20);column:code;unique"`
	Name             string                  `json:"name" gorm:"type:varchar(255);column:name"`
	WebsiteURL       string                  `json:"website_url" gorm:"type:varchar(255);column:website_url"`
	DiscordServerURL *string                 `json:"discord_server_url" gorm:"type:varchar(100);column:discord_server_url"`
	WithinVATEUD     bool                    `json:"-" gorm:"type:tinyint(1);column:within_eud;default:true"`
	Instructors      []SubdivisionInstructor `json:"instructors,omitempty" gorm:"foreignKey:SubdivisionID"`
	CreatedAt        time.Time               `json:"-" gorm:"type:timestamp;column:created_at"`
	UpdatedAt        time.Time               `json:"-" gorm:"type:timestamp;column:updated_at"`
}
