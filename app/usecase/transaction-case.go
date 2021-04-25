package usecase

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/mashingan/smapping"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/repository"
)

type TransactionCase interface {
	Insert(input dto.InsertTransactionDTO) (entity.Transaction, error)
	Fetch() ([]entity.Transaction, error)
	GetByID(id uint64) (entity.Transaction, error)
	Update(input dto.UpdateTransactionDTO) (entity.Transaction, error)
	Delete(id uint64, deleted_at time.Time) (entity.Transaction, error)
	BuyEvent(input dto.BuyEventDTO) (entity.Transaction, error)
}

type transactionCase struct {
	transactionRepository repository.TransactionRepository
	eventRepository       repository.EventRepository
}

func NewTransactionCase(transactionRepository repository.TransactionRepository, eventRepository repository.EventRepository) TransactionCase {
	return &transactionCase{
		transactionRepository: transactionRepository,
		eventRepository:       eventRepository,
	}
}

func (service *transactionCase) Insert(input dto.InsertTransactionDTO) (entity.Transaction, error) {
	transaction := entity.Transaction{}
	if err := smapping.FillStruct(&transaction, smapping.MapFields(&input)); err != nil {
		log.Println(err)
	}

	resTransaction, err := service.transactionRepository.Insert(transaction)

	if err != nil {
		log.Println(err)
	}

	return resTransaction, err
}

func (service *transactionCase) Fetch() ([]entity.Transaction, error) {
	transactions, err := service.transactionRepository.Fetch()
	if err != nil {
		log.Println(err)
	}
	return transactions, err
}

func (service *transactionCase) GetByID(id uint64) (entity.Transaction, error) {
	transaction, err := service.transactionRepository.GetByID(id)
	if err != nil {
		log.Println(err)
	}
	return transaction, err
}

func (service *transactionCase) Update(input dto.UpdateTransactionDTO) (entity.Transaction, error) {
	transaction := entity.Transaction{}

	if err := smapping.FillStruct(&transaction, smapping.MapFields(&input)); err != nil {
		log.Println(err)
	}

	resTransaction, err := service.transactionRepository.Update(transaction)
	if err != nil {
		log.Println(err)
	}

	return resTransaction, err
}

func (service *transactionCase) Delete(id uint64, deleted_at time.Time) (entity.Transaction, error) {
	transaction := entity.Transaction{}
	transaction.ID = id
	transaction.DeletedAt.Time = deleted_at

	resTransaction, err := service.transactionRepository.Delete(transaction)
	if err != nil {
		log.Println(err)
	}
	return resTransaction, err
}

func (service *transactionCase) BuyEvent(input dto.BuyEventDTO) (entity.Transaction, error) {
	if transaction, err := service.transactionRepository.GetByParticipantAndEventId(input.ParticipantId, input.EventID); err == nil {
		log.Println(transaction.StatusPayment, ":a ", reflect.TypeOf(transaction.StatusPayment))
		if transaction.StatusPayment == "passed" {
			log.Println(transaction.StatusPayment, ":b ", reflect.TypeOf(transaction.StatusPayment))

			return transaction, errors.New("aYou can't buy same ticket more than one")
		} else if transaction.StatusPayment == "failed" {
			log.Println(transaction.StatusPayment, ":c ", reflect.TypeOf(transaction.StatusPayment))

			return transaction, errors.New("Your transaction failed.")
		} else {
			log.Println(transaction.StatusPayment, ":d ", reflect.TypeOf(transaction.StatusPayment))
			return transaction, errors.New("You already checked out this ticket, complete your transaction please.")

		}
	} else {
		var transaction entity.Transaction
		event, err := service.eventRepository.GetByID(input.EventID)
		if err != nil {
			return transaction, errors.New("Event not found")
		}
		if event.Quantity == 0 {
			return transaction, errors.New("Event is out of stock")
		}

		if err := smapping.FillStruct(&transaction, smapping.MapFields(&input)); err != nil {
			log.Println(err)
		}

		resTransaction, err := service.transactionRepository.Insert(transaction)
		if err == nil {
			event.Quantity = event.Quantity - 1
			service.eventRepository.Update(event)
		}
		return resTransaction, err
	}
}
