package repository

import (
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
	tx := db.connection.Find(&users)
	return users, tx.Error
}

func (db *userRepository) GetByID(id uint64) (entity.User, error) {
	var user entity.User
	tx := db.connection.Where("id = ?", id).Take(&user)
	return user, tx.Error
}

func (db *userRepository) Update(user entity.User) (entity.User, error) {
	tx := db.connection.Save(&user)
	return user, tx.Error
}

func (db *userRepository) GetByUsername(username string) (entity.User, error) {
	var user entity.User
	tx := db.connection.Where("username = ?", username).Take(&user)
	return user, tx.Error
}

func (db *userRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	tx := db.connection.Where("email = ?", email).Take(&user)
	return user, tx.Error
}

func (db *userRepository) GetByUnique(username string, email string) (entity.User, error) {
	var user entity.User
	tx := db.connection.Where("username = ? OR email = ? ", username, email).Take(&user)
	return user, tx.Error
}

func (db *userRepository) Insert(user entity.User) (entity.User, error) {
	tx := db.connection.Create(&user)
	return user, tx.Error
}

func (db *userRepository) Delete(user entity.User) (entity.User, error) {
	tx := db.connection.Model(&user).Update("deleted_at", user.DeletedAt.Time)
	return user, tx.Error
}
