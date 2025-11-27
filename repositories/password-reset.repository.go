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
	result := r.DB.Where("user_id = ?", userID).Find(&users)

	return users, result.Error

}
