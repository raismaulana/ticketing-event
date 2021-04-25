package entity

import (
	"database/sql"
	"time"
)

type Event struct {
	ID                uint64        `gorm:"primary_key:auto_increment" json:"id"`
	CreatorId         uint64        `json:"creator_id"`
	TitleEvent        string        `gorm:"type:varchar(250)" json:"title_event"`
	LinkWebinar       string        `gorm:"type:varchar(250)" json:"link_webinar"`
	Description       string        `gorm:"type:text" json:"description"`
	Banner            string        `gorm:"type:text" json:"banner"`
	Price             float32       `gorm:"type:float(15,2)" json:"price"`
	Quantity          int32         `gorm:"type:int" json:"quantity"`
	Status            string        `gorm:"type:enum('draft','release');default:'draft'" json:"status"`
	EventStartDate    time.Time     `gorm:"type:datetime" json:"event_start_date"`
	EventEndDate      time.Time     `gorm:"type:datetime" json:"event_end_date"`
	CampaignStartDate time.Time     `gorm:"type:datetime" json:"campaign_start_date"`
	CampaignEndDate   time.Time     `gorm:"type:datetime" json:"campaign_end_date"`
	DeletedAt         sql.NullTime  `gorm:"type:timestamp null;default:null" json:"deleted_at"`
	CreatedAt         time.Time     `gorm:"<-:create;type:timestamp not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time     `gorm:"type:timestamp not null ON UPDATE CURRENT_TIMESTAMP();default:CURRENT_TIMESTAMP;" json:"Updated_at"`
	Transaction       []Transaction `gorm:"foreignKey:EventID;references:ID" json:"transaction,omitempty"`
}

func (Event) TableName() string {
	return "Event"
}
