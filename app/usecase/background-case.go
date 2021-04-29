package usecase

import (
	"errors"

	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/repository"
)

type BackgroundCase interface {
	GetPendingTransaction() ([]entity.TransactionUserEvent, error)
	GetPromotionEvent() ([]entity.Event, []entity.User, error)
}

type backgroundCase struct {
	transactionRepository repository.TransactionRepository
	eventRepository       repository.EventRepository
	userRepository        repository.UserRepository
}

func NewBackgroundCase(transactionRepository repository.TransactionRepository, eventRepository repository.EventRepository, userRepository repository.UserRepository) BackgroundCase {
	return &backgroundCase{
		transactionRepository: transactionRepository,
		eventRepository:       eventRepository,
		userRepository:        userRepository,
	}
}

func (service *backgroundCase) GetPendingTransaction() ([]entity.TransactionUserEvent, error) {
	res, err := service.transactionRepository.GetPendingTransaction()
	return res, err
}

func (service *backgroundCase) GetPromotionEvent() ([]entity.Event, []entity.User, error) {
	users, err1 := service.userRepository.GetByRole("participant")
	events, err2 := service.eventRepository.FetchAvailable()
	var err error
	if err1 != nil || err2 != nil {
		err = errors.New("No Records")
		return events, users, err
	}
	return events, users, nil
}
