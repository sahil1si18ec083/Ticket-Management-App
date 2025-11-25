package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"not null;index"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}
