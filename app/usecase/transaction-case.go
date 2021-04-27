package usecase

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"log"
	"os"
	"reflect"
	"strconv"
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
	UploadReceipt(input dto.UploadReceipt) (entity.Transaction, error)
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
	if transaction, _ := service.transactionRepository.GetByParticipantAndEventId(input.ParticipantId, input.EventID); transaction.ID != 0 {
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

func (service *transactionCase) UploadReceipt(input dto.UploadReceipt) (entity.Transaction, error) {
	transaction := entity.Transaction{}

	unbased, err := base64.StdEncoding.DecodeString(input.ImgReceipt)
	if err != nil {
		return transaction, errors.New("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		return transaction, errors.New("Bad png")
	}
	a := strconv.Itoa(int(input.ID))
	b := strconv.Itoa(int(input.ParticipantId))
	c := strconv.Itoa(int(time.Now().Unix()))
	path := "data/" + base64.StdEncoding.EncodeToString([]byte(a+b+c))
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return transaction, errors.New("Cannot open file")
	}

	png.Encode(f, im)
	log.Println(path)
	res, errRes := service.transactionRepository.UploadReceipt(input.ID, path)
	return res, errRes
}
