package models

import "gorm.io/gorm"

type PasswordResets struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TokenHash string `json:"token_hash"`
	ExpiresAt int64  `json:"expires_at"`
}
