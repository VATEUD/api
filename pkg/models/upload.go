package models

import (
	"api/utils"
	"fmt"
	"time"
)

type Upload struct {
	ID        uint      `json:"id" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	Name      string    `json:"name" gorm:"type:varchar(255);column:name"`
	Path      string    `json:"-" gorm:"type:varchar(500);column:path"`
	Type      string    `json:"type" gorm:"type:varchar(1);column:type"`
	Public    bool      `json:"-" gorm:"type:tinyint(1);column:public"`
	DownloadURL string `json:"download_url"`
	CreatedAt time.Time `json:"-" gorm:"type:timestamp;column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"type:timestamp;column:updated_at"`
}

func (upload *Upload) SetDownloadURL() {
	upload.DownloadURL = fmt.Sprintf("%s/uploads/download/%d", utils.Getenv("API_URL", ""), upload.ID)
}
