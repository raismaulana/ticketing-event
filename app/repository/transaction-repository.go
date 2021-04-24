package repository

import (
	"github.com/raismaulana/ticketing-event/app/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Fetch() ([]entity.Transaction, error)
	GetByID(id uint64) (entity.Transaction, error)
	Update(transaction entity.Transaction) (entity.Transaction, error)
	Insert(transaction entity.Transaction) (entity.Transaction, error)
	Delete(transaction entity.Transaction) (entity.Transaction, error)
}

type transactionRepository struct {
	connection *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		connection: db,
	}
}

func (db *transactionRepository) Fetch() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	tx := db.connection.Find(&transactions)
	return transactions, tx.Error
}

func (db *transactionRepository) GetByID(id uint64) (entity.Transaction, error) {
	var transaction entity.Transaction
	tx := db.connection.Where("id = ?", id).Take(&transaction)
	return transaction, tx.Error
}

func (db *transactionRepository) Update(transaction entity.Transaction) (entity.Transaction, error) {
	tx := db.connection.Save(&transaction)
	return transaction, tx.Error
}

func (db *transactionRepository) Insert(transaction entity.Transaction) (entity.Transaction, error) {
	tx := db.connection.Create(&transaction)
	return transaction, tx.Error
}

func (db *transactionRepository) Delete(transaction entity.Transaction) (entity.Transaction, error) {
	tx := db.connection.Model(&transaction).Update("deleted_at", transaction.DeletedAt.Time)
	return transaction, tx.Error
}
