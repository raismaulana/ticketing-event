package usecase

import (
	"log"
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
}

type transactionCase struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionCase(transactionRepository repository.TransactionRepository) TransactionCase {
	return &transactionCase{
		transactionRepository: transactionRepository,
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
