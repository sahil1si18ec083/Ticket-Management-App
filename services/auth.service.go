package services

import (
	"errors"
	"strconv"

	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/repositories"
	"ticket-app-gin-golang/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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
