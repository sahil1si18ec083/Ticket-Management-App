package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/repositories"
	"ticket-app-gin-golang/utils"
)

type AuthService struct {
	userRepo                   *repositories.UserRepository
	NewPasswordResetRepository *repositories.PasswordResetRepository
}

func NewAuthService(userRepo *repositories.UserRepository, passwordResetRepo *repositories.PasswordResetRepository) *AuthService {
	return &AuthService{
		userRepo:                   userRepo,
		NewPasswordResetRepository: passwordResetRepo,
	}
}

func (s *AuthService) Signup(name, email, password string) (string, error) {

	// check duplicate email
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return "", errors.New("user with this email already exists")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	// create user
	if err := s.userRepo.CreateUser(&user); err != nil {
		return "", err
	}

	// generate token
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *AuthService) Login(email, password string) (string, error) {

	// lookup user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// check password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// generate token
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *AuthService) RequestPasswordReset(email string) (string, error) {

	user, err := s.userRepo.FindByEmail(email)
	fmt.Println("bye")
	fmt.Println(user.ID)
	fmt.Println("hi")
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	token, err := utils.GenerateResetToken()
	if err != nil {
		return "", errors.New("failed to generate reset token")
	}
	hashedToken, err := utils.HashPassword(token)
	if err != nil {
		return "", errors.New("failed to hash reset token")
	}
	fmt.Println("Generated reset token:", hashedToken)

	err = s.NewPasswordResetRepository.Create(&models.PasswordResets{
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(utils.PasswordResetTokenExpiry),
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Println(err)
		return "", errors.New("failed to create password reset record")
	}

	// Send password reset email
	if err := s.sendPasswordResetEmail(email, token); err != nil {
		fmt.Println("Warning: failed to send password reset email:", err)
		// Don't fail the request if email sending fails
		// The reset token is still valid in the database
		return "", errors.New("failed to send password reset email")
	}

	return token, nil

}
func (s *AuthService) sendPasswordResetEmail(email, token string) error {
	// Placeholder for email sending logic
	// In a real application, integrate with an email service provider here
	fmt.Printf("Sending password reset email to %s with token: %s\n", email, token)
	return nil
}
