package central

import "time"

type OauthClient struct {
	ID                   uint      `json:"id" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	UserID               int       `json:"user_id" gorm:"primaryKey;type:bigint(20);unsigned;column:user_id"`
	Name                 string    `json:"name" gorm:"type:varchar(255);column:name"`
	Secret               string    `json:"-" gorm:"type:varchar(100);column:secret"`
	Provider             string    `json:"provider" gorm:"type:varchar(255);column:provider"`
	Redirect             string    `json:" redirect" gorm:"type:varchar(65535);column:redirect"`
	PersonalAccessClient string    `json:"personal_access_client" gorm:"type:boolean;column:personal_access_client"`
	PasswordClient       string    `json:"password_client" gorm:"type:boolean;column:password_client"`
	Revoked              string    `json:"revoked" gorm:"type:boolean;column:revoked"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}
