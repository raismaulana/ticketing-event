package usecase

import (
	"log"
	"time"

	"github.com/mashingan/smapping"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/repository"
)

type UserCase interface {
	Fetch() ([]entity.User, error)
	GetByID(id uint64) (entity.User, error)
	Update(input dto.UpdateUserDTO) (entity.User, error)
	Delete(id uint64, deleted_at time.Time) (entity.User, error)
	AllEventReport() ([]entity.User, error)
}

type userCase struct {
	userRepository repository.UserRepository
}

func NewUserCase(userRepository repository.UserRepository) UserCase {
	return &userCase{
		userRepository: userRepository,
	}
}

func (service *userCase) Fetch() ([]entity.User, error) {
	users, err := service.userRepository.Fetch()
	if err != nil {
		log.Println(err)
	}
	return users, err
}

func (service *userCase) GetByID(id uint64) (entity.User, error) {
	user, err := service.userRepository.GetByID(id)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

func (service *userCase) Update(input dto.UpdateUserDTO) (entity.User, error) {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(&input))
	if err != nil {
		log.Println(err)
	}
	if user.Password == "" {
		a, _ := service.userRepository.GetByID(user.ID)
		user.Password = a.Password
	}
	user.Password = helper.PasswordHash(user.Password)
	resUser, err := service.userRepository.Update(user)
	if err != nil {
		log.Println(err)
	}
	return resUser, err
}

func (service *userCase) Delete(id uint64, deleted_at time.Time) (entity.User, error) {
	user := entity.User{}
	user.ID = id
	user.DeletedAt.Time = deleted_at

	resUser, err := service.userRepository.Delete(user)
	if err != nil {
		log.Println(err)
	}
	return resUser, err
}

func (service *userCase) AllEventReport() ([]entity.User, error) {
	users, err := service.userRepository.GetAllUserJoinEvent()
	return users, err
}
