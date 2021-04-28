package repository

import (
	"log"
	"time"

	"github.com/raismaulana/ticketing-event/app/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Fetch() ([]entity.User, error)
	GetByID(id uint64) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	GetByUsername(username string) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	GetByUnique(username string, email string) (entity.User, error)
	Insert(user entity.User) (entity.User, error)
	Delete(user entity.User) (entity.User, error)
	GetAllUserJoinEvent() ([]entity.User, error)
	GetParticipant(creator_id string) ([]entity.Participant, error)
}

type userRepository struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		connection: db,
	}
}

func (db *userRepository) Fetch() ([]entity.User, error) {
	var users []entity.User
	// tx := db.connection.Find(&users)
	tx := db.connection.Raw("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL").Scan(&users)
	return users, tx.Error
}

func (db *userRepository) GetByID(id uint64) (entity.User, error) {
	var user entity.User
	// tx := db.connection.Where("id = ?", id).Take(&user)
	tx := db.connection.Raw("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND `users`.`id` = ?", id).Scan(&user)
	return user, tx.Error
}

func (db *userRepository) Update(user entity.User) (entity.User, error) {
	// tx := db.connection.Save(&user)
	tx := db.connection.Exec("UPDATE `users` SET `users`.`username` = @Username, `users`.`fullname` = @Fullname, `users`.`email` = @Email, `users`.`password` = @Password, `users`.`role` = @Role WHERE `users`.`id` = @ID", user)
	return user, tx.Error
}

func (db *userRepository) GetByUsername(username string) (entity.User, error) {
	var user entity.User
	// tx := db.connection.Where("username = ?", username).Take(&user)
	tx := db.connection.Raw("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND `users`.`username` = ?", username).Scan(&user)
	return user, tx.Error
}

func (db *userRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	// tx := db.connection.Where("email = ?", email).Take(&user)
	tx := db.connection.Raw("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND `users`.`email` = ?", email).Scan(&user)
	return user, tx.Error
}

func (db *userRepository) GetByUnique(username string, email string) (entity.User, error) {
	var user entity.User
	// tx := db.connection.Where("username = ? OR email = ? ", username, email).Take(&user)
	tx := db.connection.Raw("SELECT * FROM `users` WHERE `users`.`username` = ? AND `users`.`email` = ?", username, email).Scan(&user)
	return user, tx.Error
}

func (db *userRepository) Insert(user entity.User) (entity.User, error) {
	// tx := db.connection.Create(&user)
	tx := db.connection.Exec("INSERT INTO `users` (`username`, `fullname`, `email`, `password`, `role`) VALUES (@Username, @Fullname, @Email, @Password, @Role)", user)
	return user, tx.Error
}

func (db *userRepository) Delete(user entity.User) (entity.User, error) {
	// tx := db.connection.Model(&user).Update("deleted_at", user.DeletedAt.Time)
	tx := db.connection.Exec("UPDATE `users` SET `users`.`deleted_at` = ? WHERE `users`.`id` = ?", user.DeletedAt.Time, user.ID)
	return user, tx.Error
}

func (db *userRepository) GetAllUserJoinEvent() ([]entity.User, error) {
	var users []entity.User
	var events []entity.Event
	var transactions []entity.Transaction
	// tx := db.connection.Model(entity.User{}).Joins("JOIN transaction ON transaction.participant_id = users.id").Joins("JOIN event ON event.id = transaction.event_id").Find(&users, "`transaction`.`status_payment` = 'passed' ")
	tx := db.connection.Raw("SELECT * FROM `users` JOIN `transaction` ON `users`.`id` = `transaction`.`participant_id` JOIN `event` ON `transaction`.`event_id` = `event`.`id` WHERE `transaction`.`status_payment` = 'passed'")
	// a, _ := tx.Rows()
	// log.Println(tx.Scan(&users))
	// log.Println(tx.Scan(&events))
	// log.Println(tx.Scan(&transactions))

	log.Println(users)
	log.Println(events)
	log.Println(transactions)

	return users, tx.Error
}

func (db *userRepository) GetParticipant(creator_id string) ([]entity.Participant, error) {
	var participants []entity.Participant
	tx := db.connection.Raw("SELECT p.*, t.event_id as eid FROM `users` p JOIN transaction t on p.id = t.participant_id JOIN event e on t.event_id = e.id WHERE e.creator_id = ? AND t.status_payment = 'passed' AND e.event_end_date <= ? ORDER BY t.id", creator_id, time.Now()).Scan(&participants)
	return participants, tx.Error
}
