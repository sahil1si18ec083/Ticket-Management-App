package repositories

import (
	"ticket-app-gin-golang/models"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	DB *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{DB: db}
}

func (r *PasswordResetRepository) Create(passwordReset *models.PasswordResets) error {
	return r.DB.Create(passwordReset).Error
}
func (r *PasswordResetRepository) FindActiveByUserID(userID uint) ([]models.PasswordResets, error) {
	var users []models.PasswordResets
	// Only return unused reset tokens for the user
	result := r.DB.Where("user_id = ? AND is_used = ?", userID, false).Find(&users)

	return users, result.Error

}
func (r *PasswordResetRepository) MarkedAsUsed() error {
	return nil

}

func (r *PasswordResetRepository) Update(val *models.PasswordResets) error {
	return r.DB.Save(val).Error
}
