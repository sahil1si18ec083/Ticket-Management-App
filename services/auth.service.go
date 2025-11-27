package services

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
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

func (s *AuthService) ForgetPasswordReset(email string) (string, error) {

	user, err := s.userRepo.FindByEmail(email)
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
		IsUsed:    false,
	})
	// handle above error
	if err != nil {

		fmt.Println(err)
		return "", errors.New("failed to create password reset record")
	}

	// Send password reset email
	if err := s.sendPasswordResetEmail(email, token); err != nil {
		return "", errors.New("failed to send password reset email")

	}

	return token, nil

}
func (s *AuthService) sendPasswordResetEmail(email, token string) error {

	from := os.Getenv("FROM_EMAIL_ADDRESS")
	to := []string{email}
	subject := "Password Reset Request"
	body := fmt.Sprintf("Use the following token to reset your password: %s", token)
	mesage := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, email, subject, body)
	msg := []byte(mesage)
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	auth := smtp.PlainAuth(
		"",
		from,
		emailPassword,
		smtpHost,
	)
	smtpport := os.Getenv("SMTP_PORT")
	addrstring := fmt.Sprintf("%s:%s", smtpHost, smtpport)
	err := smtp.SendMail(
		addrstring,
		auth,
		from,
		to,
		msg,
	)
	if err != nil {
		return err
	}
	return nil

}

func (s *AuthService) ResetPassword(email, token, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("invalid email or password")
	}
	userId := user.ID
	fmt.Println(userId)
	res, err := s.NewPasswordResetRepository.FindActiveByUserID(userId)
	fmt.Println(err)
	fmt.Println(res)
	if err != nil || len(res) == 0 {
		fmt.Println(err)
		return errors.New("no active tokens found")
	}
	var matchedReset *models.PasswordResets
	for _, val := range res {
		if utils.CheckPasswordHash(token, val.TokenHash) {
			matchedReset = &val
			break
		}
	}
	if matchedReset == nil {
		return errors.New("invalid or expired reset token")
	}
	currentTime := time.Now()
	if !matchedReset.ExpiresAt.After(currentTime) {
		fmt.Println("old token ")
		return errors.New("invalid or expired reset token")
	}
	matchedReset.IsUsed = true
	err = s.NewPasswordResetRepository.Update(matchedReset)
	if err != nil {
		return errors.New("unable to  expire reset token")
	}

	newhashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	user.Password = newhashedPassword
	err = s.userRepo.DB.Save(user).Error
	if err != nil {
		return errors.New("failed to update user password")
	}

	return nil
}
