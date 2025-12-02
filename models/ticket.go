package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	UserID          uint   `json:"user_id" gorm:"not null;index"`
	User            User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Status          Status `gorm:"type:varchar(20);not null;default:'NEW'" json:"status"`
	AssignedAgentID *uint  `json:"assigned_agent_id" gorm:"index"`
	AssignedAgent   User   `gorm:"foreignKey:AssignedAgentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
