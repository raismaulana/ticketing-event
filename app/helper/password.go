package helper

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}

func PasswordVerify(hashedPwd string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPassword))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
