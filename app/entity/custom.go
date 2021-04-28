package entity

type EventReport struct {
	TotalAmount      float64 `gorm:"->" json:"total_amount"`
	Event            Event   `gorm:"embedded" json:"detail_event"`
	TotalParticipant int     `gorm:"->" json:"total_participant"`
	Participant      []User  `gorm:"-" json:"participants"`
}

type Participant struct {
	User User   `gorm:"embedded"`
	Eid  uint64 `json:"eid"`
}

type TransactionUserEvent struct {
	Email string `json:"email"`
	Link  string `json:"link"`
	Eid   uint64 `json:"eid"`
}

type ReportTransaction struct {
	Pid        uint64 `json:"pid"`
	Fullname   string `json:"fullnanme"`
	Email      string `json:"email"`
	Eid        uint64 `json:"eid"`
	TitleEvent string `json:"title_event"`
}
