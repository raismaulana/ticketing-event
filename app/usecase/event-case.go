package usecase

import (
	"log"
	"time"

	"github.com/mashingan/smapping"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/repository"
)

type EventCase interface {
	Insert(input dto.InsertEventDTO) (entity.Event, error)
	Fetch() ([]entity.Event, error)
	FetchAvailable() ([]entity.Event, error)
	GetByID(id uint64) (entity.Event, error)
	Update(input dto.UpdateEventDTO) (entity.Event, error)
	Delete(id uint64, deleted_at time.Time) (entity.Event, error)
}

type eventCase struct {
	eventRepository repository.EventRepository
}

func NewEventCase(eventRepository repository.EventRepository) EventCase {
	return &eventCase{
		eventRepository: eventRepository,
	}
}

func (service *eventCase) Insert(input dto.InsertEventDTO) (entity.Event, error) {
	event := entity.Event{}
	if err := smapping.FillStruct(&event, smapping.MapFields(&input)); err != nil {
		log.Println(err)
	}

	resEvent, err := service.eventRepository.Insert(event)

	if err != nil {
		log.Println(err)
	}

	return resEvent, err
}

func (service *eventCase) Fetch() ([]entity.Event, error) {
	events, err := service.eventRepository.Fetch()
	if err != nil {
		log.Println(err)
	}
	return events, err
}

func (service *eventCase) FetchAvailable() ([]entity.Event, error) {
	events, err := service.eventRepository.FetchAvailable()
	if err != nil {
		log.Println(err)
	}
	return events, err
}

func (service *eventCase) GetByID(id uint64) (entity.Event, error) {
	event, err := service.eventRepository.GetByID(id)
	if err != nil {
		log.Println(err)
	}
	return event, err
}

func (service *eventCase) Update(input dto.UpdateEventDTO) (entity.Event, error) {
	event := entity.Event{}
	err := smapping.FillStruct(&event, smapping.MapFields(&input))
	if err != nil {
		log.Println(err)
	}
	resEvent, err := service.eventRepository.Update(event)
	if err != nil {
		log.Println(err)
	}
	return resEvent, err
}

func (service *eventCase) Delete(id uint64, deleted_at time.Time) (entity.Event, error) {
	event := entity.Event{}
	event.ID = id
	event.DeletedAt.Time = deleted_at

	resEvent, err := service.eventRepository.Delete(event)
	if err != nil {
		log.Println(err)
	}
	return resEvent, err
}
