package repository

import (
	"log"

	"github.com/raismaulana/ticketing-event/app/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Fetch() ([]entity.Transaction, error)
	GetByID(id uint64) (entity.Transaction, error)
	Update(transaction entity.Transaction) (entity.Transaction, error)
	Insert(transaction entity.Transaction) (entity.Transaction, error)
	Delete(transaction entity.Transaction) (entity.Transaction, error)
	GetByParticipantAndEventId(participantId uint64, eventId uint64) (entity.Transaction, error)
	UploadReceipt(id uint64, path string) (entity.Transaction, error)
	VerifyPayment(id uint64, status string) (entity.Transaction, error)
	GetUserAndEventByID(id uint64) (entity.TransactionUserEvent, error)
	FetchAllUserBoughtEvent() ([]entity.ReportTransaction, error)
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
	tx := db.connection.Raw("SELECT * FROM transaction WHERE deleted_at IS NULL").Scan(&transactions)
	return transactions, tx.Error
}

func (db *transactionRepository) GetByID(id uint64) (entity.Transaction, error) {
	var transaction entity.Transaction
	tx := db.connection.Raw("SELECT * FROM transaction WHERE deleted_at IS NULL AND id = ?", id).Scan(&transaction)
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

func (db *transactionRepository) GetByParticipantAndEventId(participantId uint64, eventId uint64) (entity.Transaction, error) {
	var transaction entity.Transaction
	tx := db.connection.Raw("SELECT * FROM `transaction` WHERE `transaction`.`participant_id` = ? AND `transaction`.`event_id` = ?", participantId, eventId).Scan(&transaction)
	log.Println(tx.Error)
	log.Println(transaction)
	return transaction, tx.Error
}

func (db *transactionRepository) UploadReceipt(id uint64, path string) (entity.Transaction, error) {
	tx := db.connection.Exec("UPDATE `transaction` SET `receipt` = ? WHERE `transaction`.`id` = ?", path, id)
	return entity.Transaction{}, tx.Error
}

func (db *transactionRepository) VerifyPayment(id uint64, status string) (entity.Transaction, error) {
	tx := db.connection.Exec("UPDATE transaction SET status_payment = ? WHERE id = ?", status, id)
	return entity.Transaction{}, tx.Error
}

func (db *transactionRepository) GetUserAndEventByID(id uint64) (entity.TransactionUserEvent, error) {
	var detailTransaction entity.TransactionUserEvent
	tx := db.connection.Raw("SELECT p.email as email, e.link_webinar as link, e.id as eid FROM transaction t JOIN user p ON t.participant_id = p.id JOIN event e ON t.event_id ON e.id WHERE t.id = ?", id).Scan(&detailTransaction)
	return detailTransaction, tx.Error
}
func (db *transactionRepository) FetchAllUserBoughtEvent() ([]entity.ReportTransaction, error) {
	var reportTransaction []entity.ReportTransaction
	tx := db.connection.Raw("SELECT p.id as pid, p.fullname, p.email, e.id as eid, e.title_event FROM `users` c JOIN event e on c.id = e.creator_id JOIN transaction t on e.id = t.event_id JOIN users p on t.participant_id = p.id WHERE t.status_payment = 'passed' AND t.deleted_at IS NULL").Scan(&reportTransaction)
	return reportTransaction, tx.Error
}
