package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
	Role     Role   `gorm:"type:varchar(20);not null;default:'USER'" json:"role"`
}
