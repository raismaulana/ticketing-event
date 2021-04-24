package usecase

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/repository"
)

type AuthCase interface {
	Login(input dto.LoginDTO) interface{}
	Register(input dto.RegisterUserDTO) (entity.User, error)
	IsDuplicateUnique(username string, email string) (entity.User, error)
}

type authCase struct {
	userRepository repository.UserRepository
}

func NewAuthCase(userRepository repository.UserRepository) AuthCase {
	return &authCase{
		userRepository: userRepository,
	}
}

func (service *authCase) Login(input dto.LoginDTO) interface{} {
	resUser, err := service.userRepository.GetByUsername(input.Username)
	if err != nil {
		return false
	}
	if resUser.Username == input.Username && helper.PasswordVerify(resUser.Password, input.Password) {
		return resUser
	}
	return false
}

func (service *authCase) Register(input dto.RegisterUserDTO) (entity.User, error) {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(&input))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	user.Password = helper.PasswordHash(user.Password)
	resUser, err2 := service.userRepository.Insert(user)
	return resUser, err2
}

func (service *authCase) IsDuplicateUnique(username string, email string) (entity.User, error) {
	user, err := service.userRepository.GetByUnique(username, email)
	return user, err
}
