package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytearr := []byte(password)
	bytes, err := bcrypt.GenerateFromPassword(bytearr, bcrypt.DefaultCost)
	return string(bytes), err

}
func CheckPasswordHash(password string, hashedvalue string) bool {
	storedHashedPassword := []byte(hashedvalue)
	enteredPassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(storedHashedPassword, enteredPassword)
	if err != nil {
		return false

	}
	return true
}
