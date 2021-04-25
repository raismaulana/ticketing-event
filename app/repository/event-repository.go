package repository

import (
	"time"

	"github.com/raismaulana/ticketing-event/app/entity"
	"gorm.io/gorm"
)

type EventRepository interface {
	Fetch() ([]entity.Event, error)
	FetchAvailable() ([]entity.Event, error)
	GetByID(id uint64) (entity.Event, error)
	Update(event entity.Event) (entity.Event, error)
	Insert(event entity.Event) (entity.Event, error)
	Delete(event entity.Event) (entity.Event, error)
}

type eventRepository struct {
	connection *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		connection: db,
	}
}

func (db *eventRepository) Fetch() ([]entity.Event, error) {
	var events []entity.Event
	tx := db.connection.Find(&events)
	return events, tx.Error
}

func (db *eventRepository) FetchAvailable() ([]entity.Event, error) {
	var events []entity.Event
	tx := db.connection.Raw("SELECT * FROM event WHERE `event`.`deleted_at` IS NULL AND `event`.`status` = 'release' AND (`event`.`campaign_start_date` <= ? AND `event`.`campaign_end_date` >= ?)", time.Now(), time.Now()).Scan(&events)
	return events, tx.Error
}

func (db *eventRepository) GetByID(id uint64) (entity.Event, error) {
	var event entity.Event
	tx := db.connection.Where("id = ?", id).Take(&event)
	return event, tx.Error
}

func (db *eventRepository) Update(event entity.Event) (entity.Event, error) {
	tx := db.connection.Save(&event)
	return event, tx.Error
}

func (db *eventRepository) Insert(event entity.Event) (entity.Event, error) {
	tx := db.connection.Create(&event)
	return event, tx.Error
}

func (db *eventRepository) Delete(event entity.Event) (entity.Event, error) {
	tx := db.connection.Model(&event).Update("deleted_at", event.DeletedAt.Time)
	return event, tx.Error
}
