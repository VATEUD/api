package models

import (
	"time"
)

type DivisionExaminer struct {
	ID         uint      `json:"-" gorm:"primaryKey;type:int(11);unsigned;column:id"`
	UserID     uint      `json:"user_id" gorm:"primaryKey;type:int(11);unsigned;unique;column:user_id"`
	Name       string    `json:"name" gorm:"type:varchar(255);column:name"`
	UpTo       int       `json:"-" gorm:"type:smallint(6);column:up_to"`
	UpToString string    `json:"up_to"`
	Callsign   string    `json:"callsign" gorm:"type:varchar(255);column:callsign"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

func (examiner *DivisionExaminer) SetUpTo() {
	switch examiner.UpTo {
	case 0:
		examiner.UpToString = "S1"
		break
	case 1:
		examiner.UpToString = "S2"
		break
	case 2:
		examiner.UpToString = "S3"
		break
	case 3:
		examiner.UpToString = "C1"
		break
	default:
		examiner.UpToString = "C3"
		break
	}
}
