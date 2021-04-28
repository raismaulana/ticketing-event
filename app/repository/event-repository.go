package repository

import (
	"log"
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
	GetEventReport(creatorId string, sortf string, sorta string) ([]entity.EventReport, error)
	UpdateQuantity(eid uint64)
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
	tx := db.connection.Raw("SELECT * FROM event WHERE deleted_at IS NULL").Scan(&events)
	return events, tx.Error
}

func (db *eventRepository) FetchAvailable() ([]entity.Event, error) {
	var events []entity.Event
	tx := db.connection.Raw("SELECT * FROM event WHERE `event`.`deleted_at` IS NULL AND `event`.`status` = 'release' AND (`event`.`campaign_start_date` <= ? AND `event`.`campaign_end_date` >= ?)", time.Now(), time.Now()).Scan(&events)
	return events, tx.Error
}

func (db *eventRepository) GetByID(id uint64) (entity.Event, error) {
	var event entity.Event
	tx := db.connection.Raw("SELECT * FROM event WHERE deleted_at IS NULL AND id = ?", id).Scan(&event)
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

func (db *eventRepository) GetEventReport(creatorId string, sortf string, sorta string) ([]entity.EventReport, error) {
	var eventReport []entity.EventReport
	var stmt string

	if sorta == "ASC" {
		stmt = "SELECT e.*, SUM(t.amount) `total_amount`, COUNT(t.participant_id) `total_participant` FROM `event` e LEFT JOIN transaction t ON e.id = t.event_id LEFT JOIN users p ON t.participant_id = p.id WHERE e.event_end_date <= NOW() AND e.creator_id = " + creatorId + " GROUP BY e.id ORDER BY " + sortf + " ASC"

	} else if sorta == "DESC" {
		stmt = "SELECT e.*, SUM(t.amount) `total_amount`, COUNT(t.participant_id) `total_participant` FROM `event` e LEFT JOIN transaction t ON e.id = t.event_id LEFT JOIN users p ON t.participant_id = p.id WHERE e.event_end_date <= NOW() AND e.creator_id = " + creatorId + " GROUP BY e.id ORDER BY " + sortf + " DESC"
	}
	log.Println(stmt)
	tx := db.connection.Raw(stmt).Scan(&eventReport)
	return eventReport, tx.Error
}

func (db *eventRepository) UpdateQuantity(eid uint64) {
	db.connection.Exec("Update event SET quantity = quantity+1 WHERE id = ?", eid)
}
