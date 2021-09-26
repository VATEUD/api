package models

import "time"

type StaffDepartment struct {
	ID          uint          `json:"-" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	Title       string        `json:"name" gorm:"type:varchar(100);column:title"`
	Description string        `json:"description" gorm:"type:TEXT;column:description"`
	Email       string        `json:"email" gorm:"type:varchar(255);column:email"`
	Members     []StaffMember `json:"members" gorm:"foreignKey:DepartmentID;"`
	CreatedAt   time.Time     `json:"-" gorm:"type:timestamp;column:created_at"`
	UpdatedAt   time.Time     `json:"-" gorm:"type:timestamp;column:updated_at"`
}
