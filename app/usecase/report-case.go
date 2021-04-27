package usecase

import (
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/repository"
)

type ReportCase interface {
	FetchAllReportEvent()
	FetchAllReportEventByCreator(creator_id uint64, sortf string, sorta string) ([]entity.EventReport, error)
}

type reportCase struct {
	userRepository        repository.UserRepository
	eventRepository       repository.EventRepository
	transactionRepository repository.TransactionRepository
}

func NewReportCase(userRepository repository.UserRepository, eventRepository repository.EventRepository, transactionRepository repository.TransactionRepository) ReportCase {
	return &reportCase{
		userRepository:        userRepository,
		eventRepository:       eventRepository,
		transactionRepository: transactionRepository,
	}
}

func (service *reportCase) FetchAllReportEvent() {
}

func (service *reportCase) FetchAllReportEventByCreator(creator_id uint64, sortf string, sorta string) ([]entity.EventReport, error) {
	eventReport, err := service.eventRepository.GetEventReport(creator_id, sortf, sorta)
	participants, err2 := service.userRepository.GetParticipant(creator_id)
	if err != nil || err2 != nil {
		return []entity.EventReport{}, err
	}
	for i := range eventReport {
		j := 0
		for _, v := range participants {
			if eventReport[i].Event.ID == v.Eid {
				eventReport[i].Participant = append(eventReport[i].Participant, v.User)
				j++
			}
		}
	}
	return eventReport, nil
}
