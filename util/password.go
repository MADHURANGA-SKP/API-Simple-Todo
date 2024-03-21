package uill

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//HashPassowrd returns the bycript hash
func HashPassword(password string)(string, error){
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("faild to hash password: %w", err)
	}
	return string(HashPassword), nil
}

//checkPassword check if the provide password is correct or not
func CheckPassword(Username string,HashPassword string) error {
	err :=  bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(Username))
	
	if err != nil {
		return err
	}

	return err
}