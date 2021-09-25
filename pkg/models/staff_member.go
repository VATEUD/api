package models

import "time"

type StaffMember struct {
	ID           uint      `json:"-" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	Name         string    `json:"name" gorm:"type:varchar(255);column:name"`
	Position     string    `json:"position" gorm:"type:varchar(255);column:position"`
	Callsign     string    `json:"callsign" gorm:"type:varchar(50);column:callsign"`
	DepartmentID uint      `json:"-" gorm:"type:bigint(20);unsigned;column:department_id"`
	CreatedAt    time.Time `json:"-" gorm:"type:timestamp;column:created_at"`
	UpdatedAt    time.Time `json:"-" gorm:"type:timestamp;column:updated_at"`
}
