package repository

import (
	"errors"
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
	GetEventReport(creatorId uint64, sortf string, sorta string) ([]entity.EventReport, error)
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

func (db *eventRepository) GetEventReport(creatorId uint64, sortf string, sorta string) ([]entity.EventReport, error) {
	var eventReport []entity.EventReport
	var stmt string
	if !(sortf == "t.id" || sortf == "total_participant" || sortf == "total_amount") {
		return eventReport, errors.New("Wrong Sorting Field")
	}
	if sorta == "ASC" {
		stmt = "SELECT e.*, SUM(t.amount) `total_amount`, COUNT(t.participant_id) `total_participant` FROM `transaction` t JOIN event e on e.id = t.event_id WHERE e.creator_id = ? AND t.status_payment = 'passed' AND e.event_end_date <= ? GROUP BY t.event_id ORDER BY ? ASC"
	} else if sorta == "DESC" {
		stmt = "SELECT e.*, SUM(t.amount) `total_amount`, COUNT(t.participant_id) `total_participant` FROM `transaction` t JOIN event e on e.id = t.event_id WHERE e.creator_id = ? AND t.status_payment = 'passed' AND e.event_end_date <= ? GROUP BY t.event_id ORDER BY ? DESC"
	}
	tx := db.connection.Raw(stmt, creatorId, time.Now(), sortf).Scan(&eventReport)
	return eventReport, tx.Error
}
